from prompt_toolkit.validation import Validator, ValidationError
import os
import re

from configurator import Configuration
from configurator import PREFIX_COMMAND_HANDLER, EXIT_COMMANDS, SINGLE_WORD_COMMANDS

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

# class StatementValidator(Validator):
#     def validate(self, document):
#         text = document.text
#         if text:
#             if len(text.split()) > 1:
#                 command_word, stmt = text.split(None, 1)
#                 if command_word in OP_COMMANDS:
#                     # command is present
#                     validate_setter(text)
#                 elif command_word in SPECIAL_COMMANDS:
#                     pass
#             else:
#                 # return a parameter value or command
#                 if text in SINGLEVALUE_COMMANDS:
#                     pass
#                 elif text in CONFIG_PROPERTIES:
#                     pass
#                 else:
#                     raise ValidationError(
#                         message='{} not a valid command or configuration parameter'.format(text), 
#                         cursor_position=len(text),
#                     )

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

# Validators
def run_subvalidator(text):
    """Checks if the run statement is properly formatted

    Parameters
    ----------
    text : str
        input statement

    """
    valid_keywords = ['logger', 'threads']
    valid_logger_values = ['csv', 'sqlite']
    kwargs = dict([kwarg.split('=') for kwarg in text.split(None)[1:] if '=' in kwarg])
    for k, val in kwargs.items():
        if k not in valid_keywords:
            pos = list(re.finditer(k, text))[0].end()
            raise ValidationError(
                message='{} not a valid keyword'.format(k),
                cursor_position=pos,
            )
        else:
            if k == 'logger' and val not in valid_logger_values:
                pos = list(re.finditer(val, text))[0].end()
                raise ValidationError(
                    message='{} not a valid value for {}'.format(val, k),
                    cursor_position=pos,
                )
            elif k == 'threads' and re.search(r'\d+', val) is None:
                pos = list(re.finditer(val, text))[0].end()
                raise ValidationError(
                    message='{} not a valid value for {}'.format(val, k),
                    cursor_position=pos,
                )

def create_append_subvalidator(text):
    """Checks if the create or append statement is valid

    Parameters
    ----------
    text : str
        input statement

    """
    valid_keywords = ['intrahost_model', 'fitness_model', 'transmission_model']
    args = [arg for arg in text.split(None)[1:] if '=' not in arg]
    # Check number of arguments
    if len(args) > 2:
        raise ValidationError(
            message='Expected 2 arguments, got {}'.format(len(args)),
            cursor_position=len(text),
        )
    # Check keyword
    if args[0] not in valid_keywords:
        pos = list(re.finditer(args[0], text))[0].end()
        raise ValidationError(
            message='{} not a valid argument'.format(args[0]),
            cursor_position=pos,
        )
    # Check if model_name is a valid name
    if re.search(r'^[A-Za-z0-9\_\-\*]+$', args[1]) is None:
        raise ValidationError(
            message='{} not a a valid model name'.format(args[0]),
            cursor_position=len(text),
        )

def generate_subvalidator(text):
    """Checks if the generate statement is valid

    Parameters
    ----------
    text : str
        input statement

    """
    valid_keywords = ['pathogens', 'network', 'fitness_matrix']
    args = [arg for arg in text.split(None)[1:] if '=' not in arg]
    # Check number of arguments
    if len(args) > 1:
        raise ValidationError(
            message='Expected 1 argument, got {}'.format(len(args)),
            cursor_position=len(text),
        )
    # Check keyword
    if args[0] not in valid_keywords:
        pos = list(re.finditer(args[0], text))[0].end()
        raise ValidationError(
            message='{} not a valid argument'.format(args[0]),
            cursor_position=pos,
        )

def get_set_reset_subvalidator(text):
    """Checks if the set/get statement is valid

    Parameters
    ----------
    text : str
        input statement

    """
    kwargs = dict([kwarg.split('=') for kwarg in text.split(None)[1:] if '=' in kwarg])
    for k in kwargs.keys():
        if k not in CONFIG_PROPERTIES:
            pos = list(re.finditer(k, text))[0].end()
            raise ValidationError(
                message='{} not a valid config property'.format(k),
                cursor_position=pos,
            )

def load_save_subvalidator(text):
    """Checks if the load/save statement is valid

    Parameters
    ----------
    text : str
        input statement

    """
    valid_keywords = ['configuration', 'config']
    args = [arg for arg in text.split(None)[1:] if '=' not in arg]
    # Check number of arguments
    if len(args) > 2:
        raise ValidationError(
            message='Expected 2 arguments, got {}'.format(len(args)),
            cursor_position=len(text),
        )
    # Check keyword
    if args[0] not in valid_keywords:
        pos = list(re.finditer(args[0], text))[0].end()
        raise ValidationError(
            message='{} not a valid argument'.format(args[0]),
            cursor_position=pos,
        )
    # Check if path exists if load
    if text.split(None, 1)[0] == 'load':
        if not os.path.exists(args[1]):
            raise ValidationError(
                message='Path does not exist',
                cursor_position=len(text),
            )
        elif not os.path.isfile(args[1]):
            raise ValidationError(
                message='Path does not refer to a file',
                cursor_position=len(text),
            )
    elif text.split(None, 1)[0] == 'save':
        dirpath = os.path.dirname(args[1])
        if not os.path.exists(dirpath):
            pos = list(re.finditer(dirpath, text))[0].end()
            raise ValidationError(
                message='Directory path does not exist',
                cursor_position=pos,
            )

def todb_subvalidator(text):
    args = [arg for arg in text.split(None) if '=' not in arg]

def tocsv_subvalidator(text):
    args = [arg for arg in text.split(None) if '=' not in arg]

PREFIX_COMMAND_VALIDATOR = {
    'run': run_subvalidator,
    'create': create_append_subvalidator,
    'append': create_append_subvalidator,
    'generate': generate_subvalidator,
    'set': get_set_reset_subvalidator,
    'get': get_set_reset_subvalidator,
    'reset': get_set_reset_subvalidator,
    'load': load_save_subvalidator,
    'save': load_save_subvalidator, 
    'todb': todb_subvalidator, 
    'tocsv': tocsv_subvalidator,
}

class StatementValidator(Validator):
    def validate(self, document):
        text = document.text
        if text:
            # match with single-word commands
            if text in SINGLE_WORD_COMMANDS:
                if text in EXIT_COMANNDS:
                    pass
                elif text == 'configure':
                    pass
                elif text == 'clear':
                    pass
            # match first word
            elif text.split(None, 1) in PREFIX_COMMAND_VALIDATOR.keys():
                PREFIX_COMMAND_VALIDATOR[text.split(None, 1)](text)
