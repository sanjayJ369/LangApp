Feature: meaning

    Scenario Outline: get word meaning
        When user request the meaning of a <word>
        Then they receive it's <meaning>

        Examples:
            | word      | meaning                                   |
            | abaiser   | Ivory black; animal charcoal.             |
            | fabaceous | Having the nature of a bean; like a bean. |