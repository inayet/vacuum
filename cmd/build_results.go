package cmd

import (
	"github.com/daveshanley/vacuum/cui"
	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
	"io/ioutil"
)

func BuildResults(rulesetFlag string, specBytes []byte) (*model.RuleResultSet, *motor.RuleSetExecutionResult, error) {

	// read spec and parse
	defaultRuleSets := rulesets.BuildDefaultRuleSets()

	// default is recommended rules, based on spectral (for now anyway)
	selectedRS := defaultRuleSets.GenerateOpenAPIRecommendedRuleSet()

	// if ruleset has been supplied, lets make sure it exists, then load it in
	// and see if it's valid. If so - let's go!
	if rulesetFlag != "" {

		rsBytes, rsErr := ioutil.ReadFile(rulesetFlag)
		if rsErr != nil {
			return nil, nil, rsErr
		}
		selectedRS, rsErr = cui.BuildRuleSetFromUserSuppliedSet(rsBytes, defaultRuleSets)
		if rsErr != nil {
			return nil, nil, rsErr
		}
	}

	ruleset := motor.ApplyRulesToRuleSet(&motor.RuleSetExecution{
		RuleSet: selectedRS,
		Spec:    specBytes,
	})

	resultSet := model.NewRuleResultSet(ruleset.Results)
	resultSet.SortResultsByLineNumber()
	return resultSet, ruleset, nil
}