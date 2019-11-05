package tools

import (
	"fmt"
	"os"
	"bytes"
	"io"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"errors"
)


/******************************************************
                     apng structure
	
	Author : Joshua Pohan

	Example of usage :

	files, err := ioutil.ReadDir("./input")

	logError(err)
	test := &apng.APNGModel{}

	for _, fileInfo := range files{
		f, err := os.Open("./input/" + fileInfo.Name())
		logError(err)
		test.AppendImage(f)
		test.AppendDelay(64)
		logError(err)
	}
	test.Encode()
	test.SavePNGData("result.png")

*******************************************************/

type pngData struct{
	ihdr []byte
	idat []byte
}

type APNGModel struct{
	images []image.Image
	chunks []pngData
	delays []int
	buffer []byte
}

func (ap APNGModel) PrintPNGChunks(){
	for _, png := range ap.chunks  {
		fmt.Println("IHDR")
		fmt.Println(png.ihdr)
		fmt.Println("IDAT")
		fmt.Println(png.idat)
	}	
}


func (ap APNGModel) LogPNGChunks(){

	f, _ := os.Create("log.txt")

	for _, png := range ap.chunks {
		f.Write([]byte("IHDR\n"))
		f.Write(png.ihdr)
		f.Write([]byte("\n"))
		f.Write([]byte("IDAT\n"))
		f.Write(png.idat)
		f.Write([]byte("\n"))
	}	
}


/******************************************************
	AppendImage

	Adds file image to the apng struct, uses png package
	to decode the image

	param r The file image to add to the apng model
	return error
*******************************************************/
func (ap *APNGModel) AppendImage(r io.Reader) error {
	if curPng, err := png.Decode(r); err != nil{
		return err
	} else{
		ap.images = append(ap.images, curPng)
		return nil	
	}	
}


/******************************************************
	AppendDelay

	Adds delay in milliseconds between each image

	param delay The time delay(ms) between each images
*******************************************************/
func (ap *APNGModel) AppendDelay(delay int){
	ap.delays = append(ap.delays, delay)
}


func (ap *APNGModel) getPNGChunk(imgBuffer *bytes.Buffer) (pngData, error){
	chunk := pngData{}

	//skip png header
	imgBuffer.Next(8)

	for {
		tmp := make([]byte, 8)
		_, err := io.ReadFull(imgBuffer, tmp[:8])
		if err != nil {
			if err != io.EOF{
				return chunk, err
			} else{
				break
			}
		}

		length := binary.BigEndian.Uint32(tmp[:4])
	
		tmpVal := make([]byte, length)
		io.ReadFull(imgBuffer, tmpVal)

		switch string(tmp[4:8]){
		case "IHDR":
			chunk.ihdr = make([]byte, length)
			copy(chunk.ihdr, tmpVal)
		case "IDAT":
			chunk.idat = append(chunk.idat, tmpVal...)
			tmpVal = nil
		default:
			// do nothing, currently the tag is ignored
		}

		//skip crc
		imgBuffer.Next(4)
	}

	return chunk, nil
}


func (ap *APNGModel) appendChunk(chunk []byte, header string, toChunk *[]byte){
	chunkLe := make([]byte, 4)
	chunkTagVal := make([]byte,0, len(chunk) + 8)

	writeUint32(chunkLe, uint32(len(chunk)))

	chunkTagVal = append(chunkTagVal, []byte(header)...)
	chunkTagVal = append(chunkTagVal, chunk...)

	writeCRC32(&chunkTagVal)
	*toChunk = append(*toChunk, chunkLe...)
	*toChunk = append(*toChunk, chunkTagVal...)
}


func (ap *APNGModel) writePNGHeader(){
	ap.buffer = append(ap.buffer, 0x89,0x50,0x4E,0x47,0x0D,0x0A,0x1A,0x0A)
}


func (ap *APNGModel) appendIHDR(chunk pngData){
	ap.appendChunk(chunk.ihdr, "IHDR", &ap.buffer)
}


func (ap *APNGModel) appendacTL(img image.Image){
	tmpBuffer := []byte{}

	//number of frames in the animation
	nbFrames := make([]byte, 4)
	writeUint32(nbFrames, uint32(len(ap.images)))
	tmpBuffer = append(tmpBuffer, nbFrames...)

	//how many times the animation is looped (0 means it'll loop infinitely)
	nbLoop := make([]byte, 4)
	writeUint32(nbLoop, 0)
	tmpBuffer = append(tmpBuffer, nbLoop...)

	ap.appendChunk(tmpBuffer, "acTL", &ap.buffer)
}


func (ap *APNGModel) appendfcTL(seqNb *int, img image.Image, delay int){

	//sequence value
	fcTLValue := make([]byte, 4)
	writeUint32(fcTLValue, uint32(*seqNb))
	//width
	appendUint32(&fcTLValue, uint32(img.Bounds().Max.X - img.Bounds().Min.X))
	//height
	appendUint32(&fcTLValue, uint32(img.Bounds().Max.Y - img.Bounds().Min.Y))
	//x_offset
	appendUint32(&fcTLValue, uint32(img.Bounds().Min.X))
	//y_offset
	appendUint32(&fcTLValue, uint32(img.Bounds().Min.Y))
	//delay_num
	appendUint16(&fcTLValue, uint16(delay))
	//delay_den
	appendUint16(&fcTLValue, uint16(100))
	//dispose_op
	appendUint8(&fcTLValue, uint8(0))
	//blend_op
	appendUint8(&fcTLValue, uint8(0))
	
	ap.appendChunk(fcTLValue, "fcTL", &ap.buffer)

	//increment sequence number for animation chunk
	*seqNb++
}


func (ap *APNGModel) appendIDAT(chunk pngData){
	ap.appendChunk(chunk.idat, "IDAT", &ap.buffer)
}


func (ap *APNGModel) appendfDAT(seqNb *int, chunk pngData){

	//sequence value
	fDatValue := make([]byte, 4)
	writeUint32(fDatValue, uint32(*seqNb))

	//fdat chunk
	fDatValue = append(fDatValue, chunk.idat...)
	ap.appendChunk(fDatValue, "fdAT", &ap.buffer)

	//increment sequence number for animation chunk
	*seqNb++
}


func (ap *APNGModel) writeIENDHeader(){
	empty := make([]byte,0)
	ap.appendChunk(empty, "IEND", &ap.buffer)
}

/******************************************************
	Encode

	Encode the images previously appended into an
	apng image
*******************************************************/
func (ap *APNGModel) Encode() error {
	if len(ap.images) != len(ap.delays){
		return errors.New("Number of delays doesn't match number of images")
	}

	seqNb := 0
	for index, img := range ap.images{
		pngEnc := &png.Encoder{}
		
		pngEnc.CompressionLevel = png.BestCompression

		curImgBuffer := new(bytes.Buffer)

		if err := pngEnc.Encode(curImgBuffer, img); err != nil{
			fmt.Println(err)
			return err
		}
		if curPngChunk, err := ap.getPNGChunk(curImgBuffer); err != nil{
			return err
		} else{
			if(index == 0){
				ap.writePNGHeader()
				ap.appendIHDR(curPngChunk)
				ap.appendacTL(img)
				ap.appendfcTL(&seqNb, img, ap.delays[index])
				ap.appendIDAT(curPngChunk)
			}else{
				ap.appendfcTL(&seqNb, img, ap.delays[index])
				ap.appendfDAT(&seqNb, curPngChunk)
			}	
		}
	}
	ap.writeIENDHeader()

	return nil
}


/******************************************************
	SavePNGData

	param path    String of path of where to save the 
			      file, along with the filename
	return error
*******************************************************/
func (ap *APNGModel) SavePNGData(path string) error {

	f, _ := os.Create(path)

	_, err := f.Write(ap.buffer)
	if err != nil {
		fmt.Println(err)
	}

	f.Close()

	return err
}

/******************************************************
	WriteBytes

	param w  Writer to write the resulting encoded bytes
			
	return error
*******************************************************/
func (ap *APNGModel) WriteBytes(w io.Writer) error {

	_, err := w.Write(ap.buffer)
	if err != nil {
		fmt.Println(err)
	}

	return err
}




/******************************************************
                    byte manipulation
*******************************************************/
func appendUint8(b *[]uint8, u uint8){
	tmp := make([]byte, 1)
	writeUint8(tmp, u)
	*b = append(*b, tmp...)
}

func appendUint16(b *[]uint8, u uint16){
	tmp := make([]byte, 2)
	writeUint16(tmp, u)
	*b = append(*b, tmp...)
}

func appendUint32(b *[]uint8, u uint32){
	tmp := make([]byte, 4)
	writeUint32(tmp, u)
	*b = append(*b, tmp...)
}

func writeUint8(b []uint8, u uint8) {
	b[0] = uint8(u)
}

func writeUint16(b []uint8, u uint16) {
	b[0] = uint8(u >> 8)
	b[1] = uint8(u)
}

func writeUint32(b []uint8, u uint32) {
	b[0] = uint8(u >> 24)
	b[1] = uint8(u >> 16)
	b[2] = uint8(u >> 8)
	b[3] = uint8(u)
}

func writeCRC32(data *[]byte){
	crcBytes := make([]byte, 4)
	crc := crc32.NewIEEE()
	crc.Write(*data)
	writeUint32(crcBytes, crc.Sum32())
	*data = append(*data, crcBytes...)
}
