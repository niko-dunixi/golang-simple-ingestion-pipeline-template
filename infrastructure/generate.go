//go:build generate

//go:generate go run -tags generate .
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
)

type awsVpcModel struct {
	choices []vpc
	cursor  int

	handler vpcSelectionHandler
}

type vpc struct {
	id   string
	cidr string
	name string
}

func main() {
	log := zerolog.New(os.Stderr)
	if _, err := os.Stat("cdk.context.json"); err == nil {
		log.Info().Msg("cdk.context.json already exists. Nothing to generate here.")
		return
	}

	vpcs, err := loadVPCs()
	if err != nil {
		log.Fatal().Err(err).Msg("could not retrieve vpcs")
	}
	model := initModel(vpcs, handleSelectedVPC)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal().Err(err).Msg("could not start interactive program")
	}
}

type vpcSelectionHandler func(string) error

func loadVPCs() ([]vpc, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg)
	describedVpcResponse, err := client.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}
	results := make([]vpc, 0, len(describedVpcResponse.Vpcs))
	for _, currentVpc := range describedVpcResponse.Vpcs {
		var name string
		for _, tag := range currentVpc.Tags {
			if tag.Key != nil && *tag.Key == "Name" && tag.Value != nil {
				name = fmt.Sprintf("[%s]", *tag.Value)
			}
		}
		results = append(results, vpc{
			id:   *currentVpc.VpcId,
			cidr: *currentVpc.CidrBlock,
			name: name,
		})
	}
	return results, nil
}

func handleSelectedVPC(vpcID string) error {
	jsonBytes, err := json.Marshal(map[string]any{
		"vpc-id": vpcID,
	})
	if err != nil {
		return err
	}
	return os.WriteFile("cdk.context.json", jsonBytes, 0664)
}

func initModel(vpcs []vpc, handler vpcSelectionHandler) awsVpcModel {
	return awsVpcModel{
		choices: vpcs,
		handler: handler,
	}
}

func (m awsVpcModel) Init() tea.Cmd {
	return nil
}

func (m awsVpcModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "down", "j":
			m.cursor = (m.cursor + 1) % len(m.choices)
		case "enter", " ":
			selectedVpcID := m.choices[m.cursor].id
			if err := m.handler(selectedVpcID); err != nil {
				// FIXME: bubbletea err handling
				// What is the correct "bubbletea-way" to handle
				// this kind of error instead of a fatal?
				log.Fatalf("could not save selection (%s): %+v", selectedVpcID, err)
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m awsVpcModel) View() string {
	s := "Select target VPC:\n\n"
	for i, choice := range m.choices {
		cursor := ""
		if m.cursor == i {
			cursor = " -> "
		}
		s += fmt.Sprintf("%3s %25s %15s %s\n", cursor, choice.id, choice.cidr, choice.name)
	}
	s += "\nPress q to quit.\n"
	return s
}
