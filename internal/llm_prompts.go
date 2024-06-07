package internal

const (
	DESC_PROMPT = `Generate a description for a column in a specific table in a data warehouse.
	Criteria of a good response from you:
	- concise, 1 to 3 sentences
  - able to inform both business users and technical data analyts about the purpose and contents of the column
	- no assumptions about the data, just use business context, the table name, and the column to generate the description
	- no title, no formatting, just 1 to 3 sentences
  - avoid using the column name in the description
	- do not use tautological descriptions. Bad examples are: 'order_id' column -> "This is the id of an order"
	For example, when a table named 'orders' with a column named 'order_id', I want response like this 'The primary key of the orders table, each distinct order has a unique order_id.'
  Now, give me description for the table called %s and the column called %s.`
	TESTS_PROMPT = `Generate a list of tests that can be run on a column in a specific table in a data warehouse,
the table is called %s and the column is called %s. The tests are YAML config, there are 2 to choose from.
They have the following structure, follow this structure exactly:
  - unique
  - not_null
Return only the tests that are applicable to the column, for example, a column that is a primary key should have 
both unique and not_null tests, while a column that is a foreign key should only have the not_null test. If a 
column is potentially optional, then it should have neither test. Return only the tests that are applicable to the column.
They will be nested under a 'tests' key in a YAML file, so no need to add a title or key, just the list of tests by themselves.
  For example, a good response for a 'product_type' column in an 'orders' table would be:
  - not_null

  A good response for an 'order_id' column in an 'orders' table would be:
  - unique
  - not_null

  A good response for a 'product_sku' column in an 'orders' table would be:
  - not_null
`
)
