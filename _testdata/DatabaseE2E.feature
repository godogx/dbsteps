Feature: Database End To End

  Scenario: Successful Query
    Given all rows are deleted in table "users"

    And these rows are stored in table "users"
      | name  | email       | age | created_at           | deleted_at           |
      | Jane  | abc@aaa.com | 23  | 2021-01-01T00:00:00Z | NULL                 |
      | John  | def@bbb.de  | 33  | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | Junie | hij@ccc.ru  | 43  | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    Then only these rows are available in table "users"
      | id   | name   | email       | age | created_at           | deleted_at           |
      | $id1 | Jane   | abc@aaa.com | 23  | 2021-01-01T00:00:00Z | NULL                 |
      | $id2 | $name2 | def@bbb.de  | 33  | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | $id3 | Junie  | hij@ccc.ru  | 43  | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    And these rows are available in table "users"
      | id   | name  | email       | age | created_at           | deleted_at           |
      | 1    | Jane  | abc@aaa.com | 23  | 2021-01-01T00:00:00Z | NULL                 |
      | 2    | John  | def@bbb.de  | 33  | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | $id3 | Junie | hij@ccc.ru  | 43  | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    And no rows are available in table "orders"

    Given rows from this file are stored in table "orders"
     """
     _testdata/orders.csv
     """

    Then rows from this file are available in table "orders"
     """
     _testdata/orders.csv
     """

    Then only rows from this file are available in table "orders"
     """
     _testdata/orders.csv
     """

    Then only these rows are available in table "orders"
      | id    | amount | items | created_at | deleted_at |
      | $oid1 | $a1    | 10    | $created1  | NULL       |
      | $oid2 | $a2    | 20    | $created2  | NULL       |
      | $oid3 | $a3    | 30    | $created3  | NULL       |

    And variables are equal to values
      | $id1   | 1      |
      | $id2   | 2      |
      | $id3   | 3      |
      | $name2 | "John" |
      | $oid1  | 1      |
      | $oid2  | 2      |
      | $oid3  | 3      |
      | $a1    | 1000   |
      | $a2    | 2000   |
      | $a3    | 3000   |
