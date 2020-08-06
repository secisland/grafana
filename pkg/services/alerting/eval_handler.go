package alerting

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/metrics"
)

// DefaultEvalHandler is responsible for evaluating the alert rule.
type DefaultEvalHandler struct {
	log             log.Logger
	alertJobTimeout time.Duration
}

// NewEvalHandler is the `DefaultEvalHandler` constructor.
func NewEvalHandler() *DefaultEvalHandler {
	return &DefaultEvalHandler{
		log:             log.New("alerting.evalHandler"),
		alertJobTimeout: time.Second * 5,
	}
}

func IsExistInSlice(value string,Arr []string)bool{
    for _,v := range Arr {
        if v == value {
            return true
        }
    }
    return false
}

func IntersectSlice(arr1 []string,arr2 []string) (res []string){
    for _,v1 := range arr1 {
        for _,v2 := range arr2 {
            if v1 == v2 {
                res = append(res,v1)
            }
        }
    }
    return
}

// Eval evaluated the alert rule.
func (e *DefaultEvalHandler) Eval(context *EvalContext) {
	firing := true
	noDataFound := true
	conditionEvals := ""
	tmpEvalMatchMetric := []string{}
	tmpEvalMatches := []*EvalMatch{}

	for i := 0; i < len(context.Rule.Conditions); i++ {
		condition := context.Rule.Conditions[i]
		cr, err := condition.Eval(context)
		if err != nil {
			context.Error = err
		}

		// break if condition could not be evaluated
		if context.Error != nil {
			break
		}

		//fmt.Println("###############condition:",condition)
		//fmt.Println("###############ConditionResult:",cr)
		//for _,v := range cr.EvalMatches {
		//	fmt.Println("###############222ConditionResult-EvalMatch:",v)
		//	fmt.Println("###############222ConditionResult-EvalMatch.Metric:",v.Metric)
		//}

		if i == 0 {
			firing = cr.Firing
			noDataFound = cr.NoDataFound
			for _,v := range cr.EvalMatches {
				tmpEvalMatchMetric = append(tmpEvalMatchMetric,v.Metric)
			}
		} else {
			tmpEvalMatchMetric2 := []string{}
			for _,v := range cr.EvalMatches {
				tmpEvalMatchMetric2 = append(tmpEvalMatchMetric2,v.Metric)
			}
			tmpEvalMatchMetric = IntersectSlice(tmpEvalMatchMetric,tmpEvalMatchMetric2)
		}

		// calculating Firing based on operator
		if cr.Operator == "or" {
			firing = firing || cr.Firing
			noDataFound = noDataFound || cr.NoDataFound
		} else {
			firing = firing && cr.Firing
			if len(tmpEvalMatchMetric) ==0 {
				firing = false
			}
			noDataFound = noDataFound && cr.NoDataFound
		}

		if i > 0 {
			conditionEvals = "[" + conditionEvals + " " + strings.ToUpper(cr.Operator) + " " + strconv.FormatBool(cr.Firing) + "]"
		} else {
			conditionEvals = strconv.FormatBool(firing)
		}

		if firing {
			tmpEvalMatches = []*EvalMatch{}
			for _,v := range cr.EvalMatches {
				if IsExistInSlice(v.Metric,tmpEvalMatchMetric) {
					tmpEvalMatches = append(tmpEvalMatches,v)
					//fmt.Println("$$$$$$$$$$$$$$$$$$$:",v.Metric)
				}
			}
		} else {
			context.EvalMatches = append(context.EvalMatches, cr.EvalMatches...)
		}
	}

	if firing {
		fmt.Println("*************** alerting:")
		for _,v := range tmpEvalMatchMetric{
			fmt.Println("*************** Metric1:",v)
		}

		for _,v := range tmpEvalMatches{
			context.EvalMatches = append(context.EvalMatches,v)
			fmt.Println("*************** EvalMatcheMetric2:",v.Metric)
		}
	} else {
		fmt.Println("*************** No alerting")
	}
	
	context.ConditionEvals = conditionEvals + " = " + strconv.FormatBool(firing)
	context.Firing = firing
	context.NoDataFound = noDataFound
	context.EndTime = time.Now()

	elapsedTime := context.EndTime.Sub(context.StartTime).Nanoseconds() / int64(time.Millisecond)
	metrics.MAlertingExecutionTime.Observe(float64(elapsedTime))
}
