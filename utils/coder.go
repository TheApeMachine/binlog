package utils

import (
	"bytes"
	"context"
	"encoding/gob"

	"gitlab.itoodev.com/wrkgrp/spdg"
)

func Encode(raw interface{}) (bytes.Buffer, error) {
	var out bytes.Buffer
	enc := gob.NewEncoder(&out)
	err := enc.Encode(raw)
	return out, err
}

func Decode(raw interface{}) {
	var out bytes.Buffer
	dec := gob.NewDecoder(&out)
	dec.Decode(&raw)
}

func Wrap(ctx context.Context, out bytes.Buffer) spdg.Type {
	layer := spdg.NewProtoLayer(ctx)
	value := spdg.NewProtoValue(ctx)
	stype := spdg.NewProtoType(ctx)

	layer.Poke(out)
	value.Poke(layer)
	stype.Poke(value)

	return stype
}

func Unwrap(ctx context.Context) (spdg.Type, spdg.Value, spdg.Layer) {
	return spdg.NewProtoType(ctx), spdg.NewProtoValue(ctx), spdg.NewProtoLayer(ctx)
}

func Pushwrap(t spdg.Type, v spdg.Value, l spdg.Layer, out bytes.Buffer) spdg.Type {
	l.Poke(out)
	v.Poke(l)
	t.Poke(v)

	return t
}
