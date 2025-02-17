package mp4

import (
    "encoding/binary"
    "io"
)

// aligned(8) class TrackFragmentBaseMediaDecodeTimeBox extends FullBox(‘tfdt’, version, 0) {
// 	if (version==1) {
// 		  unsigned int(64) baseMediaDecodeTime;
// 	   } else { // version==0
// 		  unsigned int(32) baseMediaDecodeTime;
// 	   }
// 	}

type TrackFragmentBaseMediaDecodeTimeBox struct {
    Box                 *FullBox
    BaseMediaDecodeTime uint64
}

func NewTrackFragmentBaseMediaDecodeTimeBox(fragStart uint64) *TrackFragmentBaseMediaDecodeTimeBox {
    return &TrackFragmentBaseMediaDecodeTimeBox{
        Box:                 NewFullBox([4]byte{'t', 'f', 'd', 't'}, 1),
        BaseMediaDecodeTime: fragStart,
    }
}

func (tfdt *TrackFragmentBaseMediaDecodeTimeBox) Size() uint64 {
    return tfdt.Box.Size() + 8
}

func (tfdt *TrackFragmentBaseMediaDecodeTimeBox) Decode(r io.Reader) (offset int, err error) {
    if offset, err = tfdt.Box.Decode(r); err != nil {
        return
    }

    needSize := tfdt.Box.Box.Size - 12
    buf := make([]byte, needSize)
    if _, err = io.ReadFull(r, buf); err != nil {
        return 0, err
    }
    if tfdt.Box.Version == 1 {
        tfdt.BaseMediaDecodeTime = binary.BigEndian.Uint64(buf)
        offset += 8
    } else {
        tfdt.BaseMediaDecodeTime = uint64(binary.BigEndian.Uint32(buf))
        offset += 4
    }
    return
}

func (tfdt *TrackFragmentBaseMediaDecodeTimeBox) Encode() (int, []byte) {
    tfdt.Box.Box.Size = tfdt.Size()
    offset, boxdata := tfdt.Box.Encode()
    binary.BigEndian.PutUint64(boxdata[offset:], tfdt.BaseMediaDecodeTime)
    return offset + 8, boxdata
}
