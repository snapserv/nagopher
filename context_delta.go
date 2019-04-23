package nagopher

import (
	"fmt"
	"github.com/markphelps/optional"
	"reflect"
)

type deltaContext struct {
	scalarContext

	previousValue optional.Float64
}

func NewDeltaContext(name string, previousValue *float64, warningThreshold *Bounds, criticalThreshold *Bounds) Context {
	baseContext := NewScalarContext(name, warningThreshold, criticalThreshold)
	scalarContext, ok := baseContext.(*scalarContext)
	if !ok {
		panic("nagopher: Could not cast NewScalarContext() to scalarContext")
	}

	deltaContext := &deltaContext{
		scalarContext: *scalarContext,
	}

	if previousValue != nil {
		deltaContext.previousValue = optional.NewFloat64(*previousValue)
	}

	return deltaContext
}

func (c deltaContext) Evaluate(metric Metric, resource Resource) Result {
	numericMetric, ok := metric.(NumericMetric)
	if !ok {
		return NewResult(
			ResultState(StateUnknown()),
			ResultMetric(metric), ResultContext(c), ResultResource(resource),
			ResultHint(fmt.Sprintf("DeltaContext can not process metric of type [%s]", reflect.TypeOf(metric))),
		)
	}

	metricValue := numericMetric.Value()
	deltaValue := metricValue - c.previousValue.OrElse(0)
	c.previousValue.Set(metricValue)
	deltaMetric := MustNewNumericMetric(numericMetric.Name()+"_delta", deltaValue, "", nil, numericMetric.ContextName())

	emptyBounds := NewBounds()
	warningThreshold := c.warningThreshold.OrElse(emptyBounds)
	criticalThreshold := c.criticalThreshold.OrElse(emptyBounds)

	if !criticalThreshold.Match(deltaValue) {
		return NewResult(
			ResultState(StateCritical()),
			ResultMetric(deltaMetric), ResultContext(c), ResultResource(resource),
			ResultHint(criticalThreshold.ViolationHint()),
		)
	} else if !warningThreshold.Match(deltaValue) {
		return NewResult(
			ResultState(StateWarning()),
			ResultMetric(deltaMetric), ResultContext(c), ResultResource(resource),
			ResultHint(warningThreshold.ViolationHint()),
		)
	}

	return NewResult(
		ResultState(StateOk()),
		ResultMetric(deltaMetric), ResultContext(c), ResultResource(resource),
	)
}

func (c deltaContext) Performance(metric Metric, resource Resource) (OptionalPerfData, error) {
	perfData, err := NewPerfData(metric, nil, nil)
	if err != nil {
		return OptionalPerfData{}, err
	}

	return NewOptionalPerfData(perfData), nil
}
