Feature: Flashcards

  Scenario: Learner creates flashcards
    When the Learner passes some text
    Then the Learner receives flashcards from it
    And the flashcards contain words from the text
    And the Learner receives the flashcards they owns

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

  Scenario: Multiple Learner can create flashcards
    When Learner Bob passes some text
    Then Bob receives his flashcards
    When Learner Alex passes some text
    Then Alex receives his flashcards
    And Alex does not see Bobs flashcards
    And Bob does not see Alexs flashcards

  Scenario: Flashcards contain word along with it's meaning
    When the Learner passes some text
    Then they receive flashcards
    And each flashcards has meaning of the word

  Scenario: Learner can export flashcards to Anki
    When the Learner creates flashcards
    Then they can export them to Anki