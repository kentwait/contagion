"""
validators.py

Validators validate whether the input is correctly formatted
and if the values such as file paths exist.

StatementValidator handles input validation for each statement
in the interface. Within StatementValidator, commands are
validated by their respective sub-validation function.

"""
import os
import re
from prompt_toolkit.validation import Validator, ValidationError
from configurator import EXIT_COMMANDS, SINGLE_WORD_COMMANDS, CONFIG_PROPERTIES

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

def to_x_subvalidator(text):
    """Checks if the todb statement is valid

    Parameters
    ----------
    text : str
        input statement

    """
    args = [arg for arg in text.split(None) if '=' not in arg]
    # Check number of args
    if len(args) > 2:
        raise ValidationError(
            message='Expected 2 argument, got {}'.format(len(args)),
            cursor_position=len(text),
        )
    # Checks if basepath dir exists
    dirpath = os.path.dirname(args[0])
    if not os.path.exists(dirpath):
        pos = list(re.finditer(dirpath, text))[0].end()
        raise ValidationError(
            message='Directory to basepath does not exist',
            cursor_position=pos,
        )
    # Checks if outpath exists
    dirpath = os.path.dirname(args[1])
    if not os.path.exists(dirpath):
        pos = list(re.finditer(dirpath, text))[0].end()
        raise ValidationError(
            message='Directory to save path does not exist',
            cursor_position=pos,
        )

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
    'todb': to_x_subvalidator,
    'tocsv': to_x_subvalidator,
}

class StatementValidator(Validator):
    """Validates statements typed in the interface
    """
    def validate(self, document):
        text = document.text
        if text:
            # match with single-word commands
            if text in SINGLE_WORD_COMMANDS:
                if text in EXIT_COMMANDS:
                    pass
                elif text == 'configure':
                    pass
                elif text == 'clear':
                    pass
            # match first word
            elif text.split(None, 1) in PREFIX_COMMAND_VALIDATOR.keys():
                PREFIX_COMMAND_VALIDATOR[text.split(None, 1)](text)
