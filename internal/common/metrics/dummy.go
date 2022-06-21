package metrics

type NoOp struct{}

func (n NoOp) Inc(_ string, _ int) {
	// todo
}
