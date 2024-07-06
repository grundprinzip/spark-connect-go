package sql

import (
	proto "github.com/apache/spark-connect-go/v35/internal/generated"
)

// DataFrameReader supports reading data from storage and returning a data frame.
// TODO needs to implement other methods like Option(), Schema(), and also "strong typed"
// reading (e.g. Parquet(), Orc(), Csv(), etc.
type DataFrameReader interface {
	// Format specifies data format (data source type) for the underlying data, e.g. parquet.
	Format(source string) DataFrameReader
	// Load reads the underlying data and returns a data frame.
	Load(path string) (DataFrame, error)
}

// dataFrameReaderImpl is an implementation of DataFrameReader interface.
type dataFrameReaderImpl struct {
	sparkSession *sparkSessionImpl
	formatSource string
}

// NewDataframeReader creates a new DataFrameReader
func NewDataframeReader(session *sparkSessionImpl) DataFrameReader {
	return &dataFrameReaderImpl{
		sparkSession: session,
	}
}

func (w *dataFrameReaderImpl) Format(source string) DataFrameReader {
	w.formatSource = source
	return w
}

func (w *dataFrameReaderImpl) Load(path string) (DataFrame, error) {
	var format string
	if w.formatSource != "" {
		format = w.formatSource
	}
	return NewDataFrame(w.sparkSession, toRelation(path, format)), nil
}

func toRelation(path string, format string) *proto.Relation {
	return &proto.Relation{
		RelType: &proto.Relation_Read{
			Read: &proto.Read{
				ReadType: &proto.Read_DataSource_{
					DataSource: &proto.Read_DataSource{
						Format: &format,
						Paths:  []string{path},
					},
				},
			},
		},
	}
}
