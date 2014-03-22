package safeslice

type SafeSlice interface {
	Append(interface{})
	At(int)					interface{}
	Close() 				[]interface{}
	Len() 					int
	Delete(int)
	Update(int, UpdateFunc)
}

type UpdateFunc func(interface{}) interface{}

type operEnum int

const (
	addtail operEnum = iota
	at
	end
	length
	remove
	update
)

type sequenceCommand struct {
	operation	operEnum
	index		int
	value		interface{}
	result		chan<- interface{}
	data		chan<- []interface{}
	updateFunc	UpdateFunc
}

type safeSlice chan sequenceCommand

func New() SafeSlice {
	ss := make(safeSlice)
	go ss.run()
	return ss
}

func (ss safeSlice) run() {
	inner := []interface{}{}
	for cmd := range ss {
		switch cmd.operation {
			case addtail:
				inner = append(inner, cmd.value)
			case at:
				if cmd.index >= 0 && cmd.index < len(inner) {
					cmd.result <- inner[cmd.index]
				} else {
					cmd.result <- nil
				}			
			case end:
				cmd.data <- inner
				close(ss)
			case length:
				cmd.result <- len(inner)
			case remove:
				if cmd.index >= 0 && cmd.index < len(inner) {
					inner = append(inner[:cmd.index], inner[cmd.index+1:]...)
				}
			case update:
				if cmd.index >=0 && cmd.index < len(inner) {
					inner[cmd.index] = cmd.updateFunc(inner[cmd.index])
				}
		}
	}
}

func (ss safeSlice) Append(item interface{}) {
	ss <- sequenceCommand{operation: addtail, value: item}
}

func (ss safeSlice) At(ind int) interface{} {
	result := make(chan interface{})
	ss <- sequenceCommand{operation: at, index:ind, result: result}
	return <-result
}

func (ss safeSlice) Close() []interface{} {
	data := make(chan []interface{})
	ss <- sequenceCommand{operation: end, data: data}
	return <-data
}

func (ss safeSlice) Len() int {
	result := make(chan interface{})
	ss <- sequenceCommand{operation: length, result: result}
	return (<-result).(int)
}

func (ss safeSlice) Delete(ind int) {
	ss <- sequenceCommand{operation: remove, index: ind}
}

func (ss safeSlice) Update(ind int, updateFunc UpdateFunc) {
	ss <- sequenceCommand{operation: update, index: ind, updateFunc: updateFunc}
}