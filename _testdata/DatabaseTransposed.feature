Feature: Database Query Transposed

  Scenario: Successful Query (Transposed)
    Given all rows are deleted in table "my_table" of database "my_db"

    And rows from this file are stored in table "my_table" of database "my_db"
    """
    _testdata/rows.csv
    """

    And these transposed rows are stored in table "my_table" of database "my_db"
      | id         | 1                    |
      | created_at | 2021-01-01T00:00:00Z |
      | deleted_at | NULL                 |
      | foo        | foo-1                |
      | bar        | abc                  |

    Then only these transposed rows are available in table "my_table" of database "my_db"
      | id         | $id1                 | $id2                 | $id3                 |
      | foo        | $foo1                | $foo1                | foo-2                |
      | bar        | abc                  | def                  | hij                  |
      | created_at | 2021-01-01T00:00:00Z | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | deleted_at | NULL                 | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    Then only these transposed rows are available in table "my_table" of database "my_db"
      | id         | $id1                 | $id2                 | $id3                 |
      | foo        | $foo1                | $foo1                | foo-2                |
      | bar        | abc                  | def                  | hij                  |
      | created_at | 2021-01-01T00:00:00Z | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | deleted_at | NULL                 | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    And no rows are available in table "my_another_table" of database "my_db"
