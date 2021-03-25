package jsoniter

import (
	"bufio"
	"errors"
)

func (iter *Iterator) ReadStringToWriter(writer *bufio.Writer) (err error) {
	c := iter.nextToken()
	if c == '"' {
		counter := 0
		for i := iter.head; i < iter.tail; i++ {
			// require ascii string and no escape
			// for: field name, base64, number
			if iter.buf[i] == '"' {
				// fast path: reuse the underlying buffer
				_, err = writer.Write(iter.buf[iter.head:i])
				if err != nil {
					return err
				}
				iter.head = i + 1
				return writer.Flush()
			}
			counter++
		}
		_, err = writer.Write(iter.buf[iter.head:iter.tail])
		if err != nil {
			return err
		}
		iter.head = iter.tail
		counter = 0
		for iter.Error == nil {
			c := iter.readByte()
			if c == '"' {
				return writer.Flush()
			}
			counter++
			err = writer.WriteByte(c)
			if err != nil {
				return err
			}
		}
		return writer.Flush()
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	iter.ReportError("CopyStringToWriter", `expects " or n, but found `+string([]byte{c}))
	return errors.New(`CopyStringToWriter: expects " or n, but found ` + string([]byte{c}))
}
