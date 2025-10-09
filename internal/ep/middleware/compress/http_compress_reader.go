package compress

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
)

type CustomReader struct {
	SourceReader   io.ReadCloser
	CompressReader io.Reader
}

func NewCustomReader(source io.ReadCloser, encoding string) (*CustomReader, error) {
	compressReader, err := getDecompressReader(encoding, source)
	if err != nil {
		return nil, err
	}

	return &CustomReader{
		SourceReader:   source,
		CompressReader: compressReader,
	}, nil
}

func (cr *CustomReader) Read(p []byte) (n int, err error) {
	return cr.CompressReader.Read(p)
}

func (cr *CustomReader) Close() error {
	err := cr.SourceReader.Close()
	if err != nil {
		return err
	}
	if rc, ok := cr.CompressReader.(io.Closer); ok {
		return rc.Close()
	}

	return nil
}

func getDecompressReader(encoding string, reader io.Reader) (io.Reader, error) {
	switch encoding {
	case encodingGzip:
		return gzip.NewReader(reader)
	//	case encodingDeflate:
	//		return zlib.NewReader(reader)
	case encodingDeflate:
		return flate.NewReader(reader), nil
	case encodingCompress:
		return lzw.NewReader(reader, lzw.LSB, 8), nil
	case encodingBrotli:
		return brotli.NewReader(reader), nil
	}

	return nil, fmt.Errorf("unknown decompress encoding: %s", encoding)
}
