package calculator

import (
	"encoding/json"
	"testing"
)

func TestCalculator_Execute(t *testing.T) {
	calc, err := NewCalculator()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	tests := []struct {
		name       string
		params     CalculateParams
		expResult  float64
		expectErr  bool
		errMessage string
	}{
		{
			name:      "Addition",
			params:    CalculateParams{Operation: "add", OperandA: 5, OperandB: 3},
			expResult: 8,
		},
		{
			name:      "Subtraction",
			params:    CalculateParams{Operation: "subtract", OperandA: 10, OperandB: 4},
			expResult: 6,
		},
		{
			name:      "Multiplication",
			params:    CalculateParams{Operation: "multiply", OperandA: 7, OperandB: 6},
			expResult: 42,
		},
		{
			name:      "Division",
			params:    CalculateParams{Operation: "divide", OperandA: 20, OperandB: 5},
			expResult: 4,
		},
		{
			name:       "Division by zero",
			params:     CalculateParams{Operation: "divide", OperandA: 10, OperandB: 0},
			expectErr:  true,
			errMessage: "division by zero not allowed",
		},
		{
			name:       "Invalid operation",
			params:     CalculateParams{Operation: "power", OperandA: 2, OperandB: 3},
			expectErr:  true,
			errMessage: "unsupported operation: power",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal parameters
			paramsJSON, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Failed to marshal parameters: %v", err)
			}

			// Execute calculation
			resultStr, err := calc.Execute(string(paramsJSON))

			// Check for expected errors
			if tt.expectErr {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				if err.Error() != tt.errMessage {
					t.Fatalf("Expected error message '%s', got '%s'", tt.errMessage, err.Error())
				}
				return
			}

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Parse result
			var resultObj CalculateResult
			if err := json.Unmarshal([]byte(resultStr), &resultObj); err != nil {
				t.Fatalf("Failed to unmarshal result: %v", err)
			}

			// Check result
			if resultObj.Result != tt.expResult {
				t.Errorf("Expected result %f, got %f", tt.expResult, resultObj.Result)
			}
		})
	}
}

func TestCalculator_GetSchema(t *testing.T) {
	calc, err := NewCalculator()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	schema := calc.GetSchema()
	if schema.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", schema.Type)
	}

	// Check if properties exist in schema
	expectedProps := []string{"operation", "operandA", "operandB"}
	for _, prop := range expectedProps {
		if _, exists := schema.Properties[prop]; !exists {
			t.Errorf("Expected property '%s' not found in schema", prop)
		}
	}
}

func TestCalculator_GetTool(t *testing.T) {
	calc, err := NewCalculator()
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	tool := calc.GetTool()
	if tool.Name != "calculate" {
		t.Errorf("Expected tool name 'calculate', got '%s'", tool.Name)
	}

	if tool.Description == "" {
		t.Errorf("Tool description should not be empty")
	}
}