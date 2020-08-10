package conditions

import (
	"math"

	"sort"

	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/tsdb"
)

// queryReducer reduces an timeserie to a nullable float
type queryReducer struct {

	// Type is how the timeserie should be reduced.
	// Ex avg, sum, max, min, count
	Type string
}

func (s *queryReducer) Reduce(series *tsdb.TimeSeries) null.Float {
	if len(series.Points) == 0 {
		return null.FloatFromPtr(nil)
	}

	value := float64(0)
	allNull := true

	switch s.Type {
	case "avg":
		validPointsCount := 0
		for _, point := range series.Points {
			if point[0].Valid {
				value += point[0].Float64
				validPointsCount++
				allNull = false
			}
		}
		if validPointsCount > 0 {
			value = value / float64(validPointsCount)
		}
	case "sum":
		for _, point := range series.Points {
			if point[0].Valid {
				value += point[0].Float64
				allNull = false
			}
		}
	case "min":
		value = math.MaxFloat64
		for _, point := range series.Points {
			if point[0].Valid {
				allNull = false
				if value > point[0].Float64 {
					value = point[0].Float64
				}
			}
		}
	case "max":
		value = -math.MaxFloat64
		for _, point := range series.Points {
			if point[0].Valid {
				allNull = false
				if value < point[0].Float64 {
					value = point[0].Float64
				}
			}
		}
	case "count":
		value = float64(len(series.Points))
		allNull = false
	case "last":
		points := series.Points
		for i := len(points) - 1; i >= 0; i-- {
			if points[i][0].Valid {
				value = points[i][0].Float64
				allNull = false
				break
			}
		}
	case "median":
		var values []float64
		for _, v := range series.Points {
			if v[0].Valid {
				allNull = false
				values = append(values, v[0].Float64)
			}
		}
		if len(values) >= 1 {
			sort.Float64s(values)
			length := len(values)
			if length%2 == 1 {
				value = values[(length-1)/2]
			} else {
				value = (values[(length/2)-1] + values[length/2]) / 2
			}
		}
	case "diff":
		allNull, value = calculateDiff(series, allNull, value, diff)
	case "percent_diff":
		allNull, value = calculateDiff(series, allNull, value, percentDiff)
	case "count_non_null":
		for _, v := range series.Points {
			if v[0].Valid {
				value++
			}
		}

		if value > 0 {
			allNull = false
		}
	}

	if allNull {
		return null.FloatFromPtr(nil)
	}

	return null.FloatFrom(value)
}

func newSimpleReducer(t string) *queryReducer {
	return &queryReducer{Type: t}
}

func calculateDiff(series *tsdb.TimeSeries, allNull bool, value float64, fn func(float64, float64) float64) (bool, float64) {
    var (
        points = series.Points
        newest  float64
        oldest float64
        count =  0
        i      int
    )
    // get the newest point
    for i = len(points) - 1; i >= 0; i-- {
        if points[i][0].Valid {
            allNull = false
            //fmt.Println("########### Break1 newest=",points[i][0].Float64," length=",len(points), " count=",count," i=",i)
            if count == 0 {
                newest = points[i][0].Float64
                count = count + 1
                if len(points) < 4 {
                    //fmt.Println("########### break11 newest=",newest," length=",len(points), " count=",count," i=",i,points[i])
                    break
                }
            } else {
                newest = (points[i][0].Float64 + newest)/2
                //fmt.Println("########### break111 newest=",newest," length=",len(points), " count=",count," i=",i)
                break
            }
        }
    }

    count = 0

    if i >= 1 {
        // get the oldest point
        points = points[0:i]
        for i := 0; i < len(points); i++ {
            if points[i][0].Valid {
                allNull = false
                //fmt.Println("########### Break2 newest=",newest," oldest=",points[i][0].Float64," length=",len(points), " count=",count," i=",i)
                if count == 0 {
                    oldest = points[i][0].Float64
                    count = count + 1
                    if len(points) < 2 {
                        //fmt.Println("########### break22 newest=",newest," oldest=",oldest," length=",len(points), " count=",count," i=",i)
                        break
                    }
                } else {
                    oldest = (oldest + points[i][0].Float64)/2
                    //fmt.Println("########### break222 newest=",newest," oldest=",oldest," length=",len(points), " count=",count," i=",i)
                    break
                }
            }
        }
        val := fn(newest, oldest)
        value = math.Abs(val)
        //fmt.Println("*************** Diff1=",value," newest=",newest," oldest=",oldest," length=",len(points), " count=",count," i=",i)
    }
    //fmt.Println("***************  Diff2=",value," newest=",newest," oldest=",oldest," length=",len(points), " count=",count," i=",i)
    return allNull, value
}

var diff = func(newest, oldest float64) float64 {
	return newest - oldest
}

var percentDiff = func(newest, oldest float64) float64 {
	return (newest - oldest) / oldest * 100
}
