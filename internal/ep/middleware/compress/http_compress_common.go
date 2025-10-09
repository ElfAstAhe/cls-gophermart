package compress

// DefaultCompressionLevel Default compression level
const DefaultCompressionLevel = 5

// Content-Type
const (
	ContentTypeApplicationJSON = "application/json"
	ContentTypeTextHTML        = "text/html"
)

// Accept-Encoding, Content-Encoding
const (
	encodingGzip     = "gzip"     // gzip
	encodingCompress = "compress" // lzw
	encodingDeflate  = "deflate"  // zlib
	encodingBrotli   = "br"       // brotli
)
