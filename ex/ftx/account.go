package ftx

import (
	"github.com/conbanwa/wstrader/ex/ftx/structs"
	"log"
)

type Positions structs.Positions

func (client *Client) GetPositions(showAvgPrice bool) (Positions, error) {
	var positions Positions
	resp, err := client._get("positions", []byte(""))
	if err != nil {
		log.Print("Error GetPositions", err)
		return positions, err
	}
	err = _processResponse(resp, &positions)
	return positions, err
}
