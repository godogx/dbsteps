Feature: Diff Closest Row

  Scenario: Diff With Vars
    Given all rows are deleted in table "users"

    And these rows are stored in table "users"
      | name  | email       | age | created_at           | deleted_at           |
      | Jane  | abc@aaa.com | 23  | 2021-01-01T00:00:00Z | NULL                 |
      | John  | def@bbb.de  | 33  | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | Junie | hij@ccc.ru  | 43  | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    And variables are set to values
      | $id_expected | 99 |

    Then these rows are available in table "users"
      | id           | name | email      | age | created_at  | deleted_at           |
      | $id_expected | John | def@bbb.de | 33  | $created_at | 2021-01-03T00:00:00Z |

  Scenario: Diff Transposed
    Given all rows are deleted in table "users"

    And these rows are stored in table "users"
      | name  | email       | age | created_at           | deleted_at           |
      | Jane  | abc@aaa.com | 23  | 2021-01-01T00:00:00Z | NULL                 |
      | John  | def@bbb.de  | 33  | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
      | Junie | hij@ccc.ru  | 43  | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |

    And variables are set to values
      | $id_expected | 2 |

    Then these transposed rows are available in table "users"
      | id         | 1                    | $id_expected         | 3                    |
      | name       | Jane                 | John                 | Junie                |
      | email      | abc@aaa.com          | def@bbb.de           | hij@ccc.ru           |
      | age        | 23                   | 32                   | 43                   |
      | created_at | 2021-01-01T00:00:00Z | 2021-01-02T00:00:00Z | $created_at          |
      | deleted_at | NULL                 | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |
