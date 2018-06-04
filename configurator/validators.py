from prompt_toolkit.validation import Validator, ValidationError
import os
import re

from configurator import Configuration

PROMPT = 'contagion> '
EXIT_COMANNDS = ['exit', 'exit()', 'quit', 'quit()', ]
SINGLEVALUE_COMMANDS = EXIT_COMANNDS + ['configure',]
SPECIAL_COMMANDS = ['create', 'reset']
OP_COMMANDS = ['set', ]
COMMANDS = OP_COMMANDS + SPECIAL_COMMANDS + EXIT_COMANNDS
CONFIG_PROPERTIES = list(Configuration().__dict__.keys())
CREATE_FUNCTIONS = ['intrahost_model', 'fitness_model', 'transmission_model']


class DirExistsValidator(Validator):
    def validate(self, document):
        text = document.text
        dirpath = os.path.dirname(text)
        if text and not os.path.exists(dirpath):
            raise ValidationError(message='{} does not exist'.format(text))

class StatementValidator(Validator):
    def validate(self, document):
        text = document.text
        if text:
            if len(text.split()) > 1:
                command_word, stmt = text.split(None, 1)
                if command_word in OP_COMMANDS:
                    # command is present
                    validate_setter(text)
                elif command_word in SPECIAL_COMMANDS:
                    pass
            else:
                # return a parameter value or command
                if text in SINGLEVALUE_COMMANDS:
                    pass
                elif text in CONFIG_PROPERTIES:
                    pass
                else:
                    raise ValidationError(
                        message='{} not a valid command or configuration parameter'.format(text), 
                        cursor_position=len(text),
                    )

def validate_setter(text):
    command_word, stmt = text.split(None, 1)
    if command_word in OP_COMMANDS:
        kv = re.split(r'\s*\=\s*', stmt)
        if len(kv) == 1:
            l = list(re.finditer(kv[0], text))
            i = l[0].end()
            raise ValidationError(
                message='{} has no specified value'.format(text),
                cursor_position=i,
            )
        elif len(kv) == 2:
            value = kv[1].lstrip("'").lstrip('"').rstrip("'").rstrip('"')
            if not value:
                raise ValidationError(
                message='{} has no specified value'.format(text),
                cursor_position=len(text),
            )
        if len(kv) == 2 and (kv[0] not in CONFIG_PROPERTIES):
            l = list(re.finditer(kv[0], text))
            i = l[0].end()
            raise ValidationError(
                message='{} not a valid configuration parameter'.format(kv[0]),
                cursor_position=i,
            )
    elif command_word not in OP_COMMANDS:
        for i, c in enumerate(text):
            if c == ' ':
                break
        raise ValidationError(
            message='{} not a valid command'.format(command_word), 
            cursor_position=i,
        ) 
