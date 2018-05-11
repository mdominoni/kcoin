Feature: Joining network as a validator
  As a user
  I want to be able to join validators set

  Background:
    Given I have the following accounts:
      | account | password | funds |
      | A       | test     | 1000  |
      | B       | test     | 1000  |

  Scenario: Start validator
    Given I have my node running using account A
    When I start validator with 5 kcoins deposit
    And I wait for my node to be synced
    And My node is already synchronised
    And the balance of A should be around 995 kcoins
    Then I should be a validator
    And I withdraw my node from validation

  Scenario: Stop mining
    Given I have my node running using account A
    And I start validator with 5 kcoins deposit
    And I wait for my node to be synced
    And My node is already synchronised
    And the balance of A should be around 995 kcoins
    And I should be a validator
    When I withdraw my node from validation
    Then There should be 5 kcoins available to me after 5 days

   Scenario: Mining rewards: basic
    Given I have my node running using account A
    And I start validator with 5 kcoins deposit
    And I wait for my node to be synced
    And My node is already synchronised
    And the balance of A should be around 995 kcoins
    And I should be a validator

     # do some transactions
    When I unlock the account A with password 'test'
    And I transfer 10 kcoin from A to B
    And the balance of A should be around 985 kcoins

    # check if some reward was generated
    # Then the balance of B should be 20 kcoins
