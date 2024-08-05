Feature: Lemmatizer

    Scenario Outline: Lemmatizerse Word
        When <word> is lemmatized
        Then it becomes it's <root> word

        Examples:
            | word      | root   |
            | Abducting | abduct |
            | Racing    | race   |
