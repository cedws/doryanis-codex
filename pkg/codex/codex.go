package codex

import (
	"context"

	"github.com/openai/openai-go"
)

func toFloat32Slice(f []float64) []float32 {
	out := make([]float32, len(f))
	for i, v := range f {
		out[i] = float32(v)
	}
	return out
}

type EmbeddingsClient[T comparable] struct {
	openaiClient openai.Client
}

func NewEmbeddingsClient[T comparable]() EmbeddingsClient[T] {
	return EmbeddingsClient[T]{
		openaiClient: openai.NewClient(),
	}
}

type (
	EmbeddingsBatchRequest[T comparable] map[T]string
	EmbeddingsBatchResult[T comparable]  map[T][]float32
)

type batchInput[T any] struct {
	Key   T
	Value string
}

type batchInputs[T any] []batchInput[T]

func (b batchInputs[T]) Values() []string {
	var values []string

	for _, input := range b {
		values = append(values, input.Value)
	}

	return values
}

func toOrderedInputs[T comparable](r EmbeddingsBatchRequest[T]) batchInputs[T] {
	var inputs []batchInput[T]

	for k, v := range r {
		inputs = append(inputs, batchInput[T]{
			Key:   k,
			Value: v,
		})
	}

	return inputs
}

func (e *EmbeddingsClient[T]) BatchEmbed(ctx context.Context, r EmbeddingsBatchRequest[T]) (EmbeddingsBatchResult[T], error) {
	orderedInputs := toOrderedInputs(r)
	values := orderedInputs.Values()

	resp, err := e.openaiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbedding3Large,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: values,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) != len(orderedInputs) {
		panic("embedding results too short")
	}

	results := make(EmbeddingsBatchResult[T], len(resp.Data))

	for _, data := range resp.Data {
		k := orderedInputs[data.Index].Key
		results[k] = toFloat32Slice(data.Embedding)
	}

	return results, nil
}
