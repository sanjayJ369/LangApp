Feature: Exporter

    Exporter exports Learner cards

    Scenario: Exporter exports Anki cards
        Given Learner has some flashcards
        When Learner Exports to Anki
        Then Cards are Exported