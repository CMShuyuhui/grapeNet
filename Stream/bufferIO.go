// 用于网络数据传输的高性能STREAM类
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2017/7/8

package grapeStream

import "errors"

const (
	defaultSize = 2048 // 默认数据包长度
)

type BufferStream struct {
	vBuffer []byte //缓存数据

	readIndex  int64 // 读入长度
	writeIndex int64 // 写入长度
}

////////////////////////////////////////////////////////
// 新建一个Stream
func NewPacker() *BufferStream {
	return &BufferStream{
		vBuffer:    make([]byte, defaultSize),
		readIndex:  0,
		writeIndex: 0, // 预留4个字节的包长度
	}
}

func BuildPacker(body []byte) (buf *BufferStream) {
	buf = NewPacker()
	buf.WriteInt32(int32(len(body))) // 写入包长度
	buf.WriteAuto(body)
	return
}

///////////////////////////////////////////////////////
// 成员函数
// 得到有效的长度 缓冲区长度
func (b *BufferStream) Len() int64 {
	return int64(len(b.vBuffer))
}

// 读到结尾了
func (b *BufferStream) EndOf() bool {
	return (b.Len() - b.readIndex) <= 0
}

// 可读取的剩余长度
func (b *BufferStream) Available() int64 {
	return (b.writeIndex - b.readIndex)
}

// 返回缓冲区全部数据
func (b *BufferStream) Bytes() []byte {
	return b.vBuffer
}

// 扩容行为
func (b *BufferStream) Resize(newSize int64) error {
	if newSize < b.Len()+defaultSize {
		return errors.New("NewSize is too short...")
	}

	tmpData := make([]byte, newSize)   // 保存构建一个新的BUFFER
	copy(tmpData[:b.Len()], b.vBuffer) // copy旧的数据进
	b.vBuffer = tmpData                // 新的缓冲区大小

	return nil
}

func (b *BufferStream) Slice(startIndex, endIndex int) []byte {
	if int64(startIndex) >= b.Len() || startIndex == endIndex {
		return b.vBuffer
	}

	if int64(endIndex) > b.Len() {
		return b.vBuffer[startIndex:]
	}

	return b.vBuffer[startIndex:endIndex]
}

// 跳过指定的字节数量
// 数据自动前移
func (b *BufferStream) Shift(size int) error {
	if int64(size) > b.writeIndex {
		return errors.New("Shift Size is too big...")
	}
	copy(b.vBuffer, b.vBuffer[size:]) // 后面的数据返回到前面
	b.readIndex -= int64(size)
	b.writeIndex -= int64(size)

	return nil
}

// 相当于不去读取而是改变pos位置
func (b *BufferStream) Skip(size int) error {
	if int64(size) > b.writeIndex {
		return errors.New("Skip Size is too big...")
	}

	b.readIndex += int64(size) // 跳过部分字节
	return nil
}

// 从开始跳到指定位置
func (b *BufferStream) Seek(pos int) error {
	if int64(pos) > b.writeIndex {
		return errors.New("Seek Pos OverFlow...")
	}
	if pos < 0 {
		return errors.New("Unknow Seek Pos...")
	}

	b.readIndex = int64(pos)
	return nil
}

// 获得当前的读位置
func (b *BufferStream) ReadPos() int64 {
	return b.readIndex
}

// 获取当前写入位置
func (b *BufferStream) WritePos() int64 {
	return b.writeIndex
}

// 读行为
// 取出数据 但不计数
func (b *BufferStream) PeekBytes(size int) []byte {
	if int64(size) > b.Available() {
		return []byte{}
	}
	return b.vBuffer[b.readIndex : b.readIndex+int64(size)]
}

func (b *BufferStream) Peek16() uint16 {
	return BTUint16(b.PeekBytes(2))
}

func (b *BufferStream) Peek32() uint32 {
	return BTUint32(b.PeekBytes(4))
}

func (b *BufferStream) GetBytes(size int) (r []byte) {
	r = b.PeekBytes(size)
	if len(r) == 0 {
		return
	}

	b.Skip(size)
	return
}

func (b *BufferStream) GetUint8() uint8 {
	return BTUint8(b.GetBytes(1))
}

func (b *BufferStream) GetInt8() int8 {
	return BTInt8(b.GetBytes(1))
}

func (b *BufferStream) GetUint16() uint16 {
	return BTUint16(b.GetBytes(2))
}

func (b *BufferStream) GetInt16() int16 {
	return BTInt16(b.GetBytes(2))
}

func (b *BufferStream) GetUint32() uint32 {
	return BTUint32(b.GetBytes(4))
}

func (b *BufferStream) GetInt32() int32 {
	return BTInt32(b.GetBytes(4))
}

func (b *BufferStream) GetUint64() uint64 {
	return BTUint64(b.GetBytes(8))
}

func (b *BufferStream) GetInt64() int64 {
	return BTInt64(b.GetBytes(8))
}

func (b *BufferStream) GetFloat32() float32 {
	return BTFloat32(b.GetBytes(4))
}

func (b *BufferStream) GetFloat64() float64 {
	return BTFloat64(b.GetBytes(8))
}

func (b *BufferStream) GetString(size int) string {
	return string(b.GetBytes(size))
}

/////////////////////////////////
// 写行为
func (b *BufferStream) WriteAuto(buf []byte) int {
	return b.Write(buf, len(buf))
}

func (b *BufferStream) Write(buf []byte, wlen int) int {
	endPos := b.writeIndex + int64(wlen)
	if endPos > b.Len() {
		err := b.Resize(endPos + defaultSize) // 扩容
		if err != nil {
			return -1
		}
	}

	copy(b.vBuffer[b.writeIndex:endPos], buf)
	b.writeIndex += int64(wlen)

	return wlen
}

func (b *BufferStream) WriteUInt8(v uint8) {
	b.WriteAuto(U8TBytes(v))
}

func (b *BufferStream) WriteInt8(v int8) {
	b.WriteAuto(I8TBytes(v))
}

func (b *BufferStream) WriteUInt16(v uint16) {
	b.WriteAuto(U16TBytes(v))
}

func (b *BufferStream) WriteInt16(v int16) {
	b.WriteAuto(I16TBytes(v))
}

func (b *BufferStream) WriteUInt32(v uint32) {
	b.WriteAuto(U32TBytes(v))
}

func (b *BufferStream) WriteInt32(v int32) {
	b.WriteAuto(I32TBytes(v))
}

func (b *BufferStream) WriteUInt64(v uint64) {
	b.WriteAuto(U64TBytes(v))
}

func (b *BufferStream) WriteInt64(v int64) {
	b.WriteAuto(I64TBytes(v))
}

func (b *BufferStream) WriteFloat32(v float32) {
	b.WriteAuto(F32TBytes(v))
}

func (b *BufferStream) WriteFloat64(v float64) {
	b.WriteAuto(F64TBytes(v))
}

func (b *BufferStream) WriteString(v string) {
	b.WriteAuto([]byte(v))
}

/////////////////////////////
// 修改指定位置数据
func (b *BufferStream) ChangeAuto(pos int, buf []byte) {
	wlen := len(buf)
	if int64(pos+wlen) > b.Len() {
		return
	}

	for i := pos; i < pos+wlen; i++ {
		b.vBuffer[i] = buf[i-pos]
	}
}

func (b *BufferStream) ChangeUInt8(pos int, v uint8) {
	b.ChangeAuto(pos, U8TBytes(v))
}

func (b *BufferStream) ChangeInt8(pos int, v int8) {
	b.ChangeAuto(pos, I8TBytes(v))
}

func (b *BufferStream) ChangeUInt16(pos int, v uint16) {
	b.ChangeAuto(pos, U16TBytes(v))
}

func (b *BufferStream) ChangeInt16(pos int, v int16) {
	b.ChangeAuto(pos, I16TBytes(v))
}

func (b *BufferStream) ChangeUInt32(pos int, v uint32) {
	b.ChangeAuto(pos, U32TBytes(v))
}

func (b *BufferStream) ChangeInt32(pos int, v int32) {
	b.ChangeAuto(pos, I32TBytes(v))
}

func (b *BufferStream) ChangeUInt64(pos int, v uint64) {
	b.ChangeAuto(pos, U64TBytes(v))
}

func (b *BufferStream) ChangeInt64(pos int, v int64) {
	b.ChangeAuto(pos, I64TBytes(v))
}

func (b *BufferStream) ChangeFloat32(pos int, v float32) {
	b.ChangeAuto(pos, F32TBytes(v))
}

func (b *BufferStream) ChangeFloat64(pos int, v float64) {
	b.ChangeAuto(pos, F64TBytes(v))
}

func (b *BufferStream) ChangeString(pos int, v string) {
	b.ChangeAuto(pos, []byte(v))
}

// 通用解包函数
// 默认协议的首部4个字节为包长度，并返回一个仅有该数据内容的STREAM
// |len 4byte|body or header|
// Unpack后会自动shift
func (b *BufferStream) Unpack(shift bool) (buf *BufferStream, err error) {
	buf = nil
	err = errors.New("Pack Unready...")
	if b.Available() < 4 {
		return // 包还没准备好
	}
	// 取出剩余长度
	len := int64(b.Peek32())
	if b.Available() < (len + 4) {
		return // 包还没准备好
	}

	b.Skip(4)                               // 跳过数据包长度
	buf = BuildPacker(b.GetBytes(int(len))) // 读取body长度
	err = nil

	if shift {
		b.Shift(int(4 + len)) // 跳过指定长度
	}

	return
}

// 通用打包体系
// |len 4byte|body or header|
func (b *BufferStream) Packer() (buf []byte, err error) {
	buf = b.Bytes()
	err = errors.New("No Data Need Package...")
	if b.writeIndex < 4 {
		return
	}

	b.ChangeUInt32(0, uint32(b.writeIndex)) // 改变长
	buf = b.Slice(0, int(b.writeIndex))
	err = nil
	return
}