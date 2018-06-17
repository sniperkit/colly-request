package httpclient

import (
	codecs "github.com/sniperkit/codecs/pkg"
	codecs "github.com/sniperkit/codecs/plugin/service"
)

var codecService services.CodecService

// Codecs returns the
// codecs "github.com/sniperkit/codecs/pkg/services".CodecService currently in use
// by this library.
func Codecs() services.CodecService {
	if codecService == nil {
		codecService = services.NewWebCodecService()
	}
	return codecService
}

// SetCodecs can be used to change the
// codecs "github.com/sniperkit/codecs/pkg/services".CodecService used by this
// library.
func SetCodecs(newService services.CodecService) {
	codecService = newService
}

// AddCodec adds a codecs "github.com/sniperkit/codecs/pkg".Codec to the
// codecs "github.com/sniperkit/codecs/pkg/services".CodecService currently in use
// by this library.
func AddCodec(codec codecs.Codec) {
	Codecs().AddCodec(codec)
}
