Feature: No locking for different tables

  Scenario: Table 1
    Given I sleep
    Given I should not be blocked for "db1::t1"
    Given there are no rows in table "t1" of database "db1"
    And I sleep

  Scenario: Table 1 again
    Given I sleep
    Given I should not be blocked for "db1::t1"
    Given there are no rows in table "t1" of database "db1"
    And I sleep

  Scenario: Table 3
    Given I sleep
    Given I should not be blocked for "db3::t3"
    Given there are no rows in table "t3" of database "db3"
    And I sleep
