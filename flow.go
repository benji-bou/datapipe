package datapipe

type Flow[I any, O any] struct {
	outputs []Outputable[O]
	input   Inputable[I]
}

// func NewFlow[I any, O any](i ...Inputable[I]) *Flow[I, O] {
// 	return &Flow[I, O]{outputs: []Outputable[O]{}, input: Merge(i), pipes: []Pipeable[any, any]{}}
// }

// func (f *Flow[I, O]) AddOutput(o ...Outputable[O]) *Flow[I, O] {
// 	f.outputs = append(f.outputs, o...)
// 	return f
// }

// func (f *Flow[I, O]) AddPipe(i ...Inputable[I]) *Flow[I, O] {
// 	f.inputs = append(f.inputs, i...)
// 	return f
// }

// func (f *Flow[I, O]) Run() <-chan error {
// 	input := Merge(f.inputs...)
// 	if len(f.pipes) > 0 {
// 		currentI = f.pipes[0].Pipe(input)
// 		for i := 1; i < len(f.pipes); i++ {
// 			cu
// 		}

// 	}
// }
