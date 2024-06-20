Feature: Flashcards

  Scenario: Learner creates flashcards
    When the Learner passes some text
    Then the Learner receives flashcards from it

  Scenario: Learner gets flashcards without duplicates
   When the Learner passes some text with repeated words
   Then the Learner receives flashcards without duplicates
