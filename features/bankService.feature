Feature: Bank Service

  # Technically, these are more Unit Tests than Integration Tests, at least in current implementation
  Scenario: Create Bank User
    Given a user provides their first name "example", last name "example", social security number "8675309-1542", and balance 150
    Then the user is created
    And has a balance of 150

  Scenario: Add User to Bank Users
    Given a user
    When the user is added to the bank users
    Then the user is in the bank users

  # You have to carefully construct the steps; i.e. if you expect some piece of data to be added into the step, you should
  # make it an explicit step for verification
  Scenario: Remove User from the Bank
    Given a user
    When the user is added to the bank users
    When the user is removed from the bank users
    Then the user is not in the bank users

#  I think is meant for Integration testing because of how I change the service being exposed
  Scenario: Transfer Money Between Two Users
    Given a user with a balance of 100
    And another user with a balance of 200
    When the first user transfers 50 to the second user
    Then the first user has a balance of 50
    And the second user has a balance of 250