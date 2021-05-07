package main

// Matrix is a new type defined by a double slice of float64.
type Matrix [][]float64

// NewMatrix creates a rows x cols matrix
func NewMatrix(rows, columns int) Matrix {
	matrix := make([][]float64, rows, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, columns, columns)
	}
	return matrix
}

// New3x3IdentityMatrix returns a 3x3 copy of that matrix.
func New3x3IdentityMatrix() Matrix {
	return Matrix(
		[][]float64{
			[]float64{1, 0, 0},
			[]float64{0, 1, 0},
			[]float64{0, 0, 1},
		},
	)
}

// MultiplyMatrixByTuple returns the multiplication of a Matrix by a Tuple.
func (matrix Matrix) MultiplyMatrixByTuple(vector *Vector3D) *Vector3D {
	tupleAsMatrix := []float64{vector.x, vector.y, vector.z}
	newTup := &Vector3D{
		matrix.dotProducOfMatricesRowColumn(matrix.Row(0), tupleAsMatrix),
		matrix.dotProducOfMatricesRowColumn(matrix.Row(1), tupleAsMatrix),
		matrix.dotProducOfMatricesRowColumn(matrix.Row(2), tupleAsMatrix),
	}

	return newTup
}

// dotProducOfMatricesRowColumn computes the dot product of a row-column combination between the two matrices.
//
// A[i] * B[i] + A[i + 1] * B[i + 1] ...
func (matrix Matrix) dotProducOfMatricesRowColumn(A, B []float64) float64 {

	length := int(min(float64(len(A)), float64(len(B))))
	total := 0.0
	for i := 0; i < length; i++ {
		total += A[i] * B[i]
	}
	return total
}

// Row returns the slice from the elements of the entire row from the current matrix.
func (matrix Matrix) Row(r int) []float64 {
	return matrix[r]
}
