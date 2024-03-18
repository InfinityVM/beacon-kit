// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package signing

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/primitives"
)

type Version [VersionLength]byte

// ForkData is the fork data used for signing.
type ForkData struct {
	CurrentVersion Version `ssz-size:"4"`
	ChainID        []byte  `ssz-max:"50"`
}

// computeForkDataRoot computes the root of the fork data.
func computeForkDataRoot(
	currentVersion Version,
	chainID string,
) (primitives.HashRoot, error) {
	forkData := ForkData{
		CurrentVersion: currentVersion,
		ChainID:        []byte(chainID),
	}
	return forkData.HashTreeRoot()
}

// VersionFromUint returns a Version from a uint32.
func VersionFromUint32(version uint32) Version {
	versionBz := Version{}
	binary.LittleEndian.PutUint32(versionBz[:], version)
	return versionBz
}