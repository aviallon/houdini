package lists_test

import (
	"testing"

	"code.houdinigraphql.com/packages/houdini-core/config"
	"code.houdinigraphql.com/packages/houdini-core/plugin"
	"code.houdinigraphql.com/plugins/tests"
)

func TestInsertOperationInput(t *testing.T) {
	tests.RunTable(t, tests.Table[config.PluginConfig, *plugin.HoudiniCore]{
		Schema: `
			input UserFilterInput {
				name: String
			}

			type Query {
				users(filter: UserFilterInput, limit: Int, offset: Int): [User!]!
			}

			type User {
				id: ID!
				firstName: String!
				field(filter: String): String
			}
		`,
		Tests: []tests.Test[config.PluginConfig]{
			{
				Name: "Operation fragments",
				Pass: true,
				Input: []string{
					`
						query AllUsers {
							users @list(name: "All_Users") {
								firstName
							}
						}
					`,
				},
				Expected: []tests.ExpectedDocument{
					tests.ExpectedDoc(`
						query AllUsers {
							users @list(name: "All_Users") {
								firstName
								__typename
								id
							}
						}
					`),
					tests.ExpectedDoc(`
						fragment All_Users_insert on User {
							firstName
							id
							__typename
						}
					`),
					tests.ExpectedDoc(`
						fragment All_Users_remove on User {
							__typename
							id
						}
					`),
					tests.ExpectedDoc(`
						fragment All_Users_toggle on User {
							firstName
							id
							__typename
						}
					`),
				},
			},
			{
				Name: "Operation fragments with variables in selection",
				Pass: true,
				Input: []string{
					`
          query AllUsers($filter: String) {
							users @list(name: "All_Users") {
                field(filter: $filter)
							}
						}
					`,
				},
				Expected: []tests.ExpectedDocument{
					tests.ExpectedDoc(`
						fragment All_Users_insert on User {
              field(filter: $filter)
							id
							__typename
						}
					`).WithVariables(tests.ExpectedOperationVariable{
						Name: "filter", Type: "String",
					}),
				},
			},
			{
				Name: "Operation fragments from @paginate",
				Pass: true,
				Input: []string{
					`
						query AllUsers {
							users(limit: 10) @paginate(name: "All_Users") {
								firstName
							}
						}
					`,
				},
				Expected: []tests.ExpectedDocument{
					tests.ExpectedDoc(`
						query AllUsers($limit: Int = 10, $offset: Int) @dedupe(match: Variables) {
							users(limit: $limit, offset: $offset) @paginate(name: "All_Users") {
								firstName
								__typename
								id
							}
						}
					`),
					tests.ExpectedDoc(`
						fragment All_Users_insert on User {
							firstName
							id
							__typename
						}
					`),
					tests.ExpectedDoc(`
						fragment All_Users_remove on User {
							id
							__typename
						}
					`),
					tests.ExpectedDoc(`
						fragment All_Users_toggle on User {
							firstName
							id
							__typename
						}
					`),
				},
			},
			{
				Name: "@paginate with object input",
				Pass: true,
				Input: []string{
					`
						query FilteredUsers {
							users(filter: { name: "Seppe" }, limit: 10) @paginate(name: "Filtered_Users") {
								firstName
							}
						}
					`,
				},
				Expected: []tests.ExpectedDocument{
					tests.ExpectedDoc(`
						query FilteredUsers($limit: Int = 10, $offset: Int) @dedupe(match: Variables) {
							users(filter: { name: "Seppe" }, limit: $limit, offset: $offset) @paginate(name: "Filtered_Users") {
								firstName
								__typename
								id
							}
						}
					`),
					tests.ExpectedDoc(`
						fragment Filtered_Users_insert on User {
							firstName
							id
							__typename
						}
					`),
					tests.ExpectedDoc(`
						fragment Filtered_Users_remove on User {
							id
							__typename
						}
					`),
					tests.ExpectedDoc(`
						fragment Filtered_Users_toggle on User {
							firstName
							id
							__typename
						}
					`),
				},
			},
		},
	})
}
