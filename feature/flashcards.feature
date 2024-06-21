Feature: Flashcards

  Scenario: Learner creates flashcards
    When the Learner passes some text
    Then the Learner receives flashcards from it

  Scenario: Learner gets flashcards without duplicates
    When the Learner passes some text with repeated words
    Then the Learner receives flashcards without duplicates

  Scenario: Learner does not get new flashcards from the same text
    When the Learner passes some text
    Then the Learner receives flashcards from it
    When the Learner passes the same text again
    Then the Learner does not receive new flashcards

  Scenario Outline: Learner can learn flashcards
    When the Learner receives a flashcard
    Then the Learner can <guess> the meaning of it
    And the flashcard becomes <memorized>
    Examples:
      | <guess> | <memorized> |
      | right   | yes         |
      | wrong   | no          |
