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

  Scenario: Multiple Learner can create flashcards
    When Learner Sanjay creates flashcards
    When Learner Dima creates flashcards
    And Dima does not see Sanjay flashcards
    And Sanjay does not see Dima flashcards

  Scenario: Flashcards contain word along with it's meaning
    When the Learner passes some text
    Then they receive flashcards
    And each flashcards has meaning of the word

  Scenario: Learner can export flashcards to Anki
    When the Learner creates flashcards
    Then they can export them to Anki