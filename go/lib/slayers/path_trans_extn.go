// Copyright 2020 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slayers

import (
    "github.com/google/gopacket"

    "github.com/scionproto/scion/go/lib/addr"
    "github.com/scionproto/scion/go/lib/common"
    "github.com/scionproto/scion/go/lib/errors"
    "github.com/scionproto/scion/go/lib/snet"
    "github.com/scionproto/scion/go/lib/spath"
)


var (
    // not all the neccessary fields were in the E2E extension
    errExtnPathTransIncomplete = serrors.New("Incomplete Path Transport Extension")
    errExtnPathTransInexistent = serrors.New("No Path Transport Extension Found"
)


type ExtnPathTrans struct{
    SrcIA   addr.IA
    SrcHost addr.HostAddr
    Path    *spath.Path
}

func NewExtnPathTransFromLayer(extn *EndToEndExtn, c *snet.Conn) (*ExtnPathTrans, error){
    var extn ExtnPathTrans
    err = ext.DecodeFromLayer(extension)
    if err != nil {
        return nil, err
    }
    return &extn, nil


}

func (o *ExtnPathTrans) DecodeFromLayer(extn *EndToEndExtn) error{
    existance bool = false
    for option := range extn.Options {
        if option.OptType == OptTypePathTransExtn {
            // if option.OptDataLen < 8{
            //     // not enough data for a full Path Transport Extension
            //     return errExtnPathTransIncomplete
            // }
            existance = true
            // copy data to extension
            offset := 0
            srcHostType HostAddrType = option.OptData[offset]
            offset += 1

            // parse the new address
            l, err := o.parseAddr(option.OptData[offset:],srcHostType)
            if err != nil{
                return nil
            }
            offset += l
            // currently no padding after SrcAddr
            o.parsePath(data[offset:])
            break

        }else{
            continue
        }

    }
    if !existance{
        // we didn't have a Path Transport Extension in the packet
        return errExtnPathTransInexistent
    }
    return nil
}

func (o *ExtnPathTrans) Write(location []byte) error{
    offset := 0

    hostType := addr.HostTypeNone
    if o.SrcHost != nil{
        // we have host addr
        hostType = o.SrcHost.Type()
    }
    location[offset] = uint8(hostType)
    offset += 1

    o.SrcIA.Write(location[offset:])
    offset += addr.IABytes

    if o.SrcHost != nil{
        srcHost := o.SrcHost.Pack()
        copy(location[offset:],srcHost)
        offset += len(srcHost)
    }else{
        offset += addr.HostLenNone
    }

    // currently no padding added after SrcAddr

    if o.Path != nil {
        copy(location[offset:],o.Path.Raw)
    }
    return nil
}

func (o *ExtnPathTrans) Len() int{
    var hostLen uint8
    if o.SrcHost != nil{
        hostLen,_ = addr.HostLen(o.SrcHost.Type())
    }

    var pathLen int
    if o.Path != nil{
        pathLen = len(o.Path.Raw)
    }
    return len(addr.HostAddrType) + addr.IABytes + int(hostLen) + pathLen
}

func (o *ExtnPathTrans) Copy() 

func (o *ExtnPathTrans) parseAddr(data []byte, srcType HostAddrType) (int,error){
    // length of src host address
    srcL := srcType.Size()

    // total length of SCION address
    totLen := addr.IABytes + srcL
    if data < totLen{
        return 0,errExtnPathTransIncomplete
    }

    // parse IA and Host addr from data
    ia := addr.IAFromRaw(data)
    host, err := addr.HostFromRaw(data[addr.IABytes:],srcType)
    if err != nil{
        return 0,err
    }

    o.SrcIA = ia
    o.SrcHost = host

    return totLen,nil
}

func (o *ExtnPathTrans) parsePath(data []byte) {
    // make empty or SCION path depending on size
    if len(data) > 0{
        o.Path.Raw = data[:]
        o.Path.Type = slayer.PathTypeSCION
    }
    else{
        o.Path.Raw = []
        o.Path.Type = slayer.PathTypeEmpty
    }
    return
}




