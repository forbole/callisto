/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"
)

// SaveIssues - method that create issues to the "vipcoin_chain_banking_issue" table
func (r Repository) SaveIssues(issues ...*bankingtypes.Issue) error {
	return nil
}

// GetIssues - method that get issues from the "vipcoin_chain_banking_issue" table
func (r Repository) GetIssues(filter filter.Filter) ([]*bankingtypes.Issue, error) {
	return []*bankingtypes.Issue{}, nil
}

// UpdateIssues - method that update the issue in the "vipcoin_chain_banking_issue" table
func (r Repository) UpdateIssues(issues ...*bankingtypes.Issue) error {
	return nil
}
