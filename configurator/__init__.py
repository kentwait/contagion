"""
The configurator module contains classes and functions related to
running the command-prompt interface in contagion.py.

The module is divided into three main files:
- validators.py
- handlers.py
- configuration.py

Validation of input statements are performed by StatementValidator
an instance of the Validator class in prompt_toolkit. Depending on
the command, StatementValidator calls specific sub-validators in
validators.py to handle validating a particular command for
formatting and value.

Commands are executed its respective handler function in handlers.py.
The handler function may prompt for more input and validate the given
data as well.

Finally, configuration.py holds the Configuration class. The Configuration
class holds the settings and parameters to create and run the simulation.
When ready, the Configuration class can write the settings into a TOML
file that can be read by Contagion. It can also import settings from an
existing TOML file.
"""
from configurator.configuration import Configuration
from configurator import handlers

# Lists the properties in the Configuration object.
# Used to validate keywords against.
CONFIG_PROPERTIES = list(Configuration().__dict__.keys())

# Commands that breaks the input loop.
EXIT_COMMANDS = ['exit', 'exit()', 'quit', 'quit()', 'q']

# Other single-word commands aside from exit.
SINGLE_WORD_COMMANDS = EXIT_COMMANDS + ['configure', 'clear']

# Map of command handlers called depending on the prefix command used.
PREFIX_COMMAND_HANDLER = {
    'run': handlers.run_handler,
    'create': handlers.create_handler,
    'append': handlers.append_handler,
    'generate': handlers.generate_handler,
    'set': handlers.set_handler,
    'get': handlers.get_handler,
    'reset': handlers.reset_handler,
    'load': handlers.load_handler,
    'save': handlers.save_handler, 
    'todb': handlers.todb_handler, 
    'tocsv': handlers.tocsv_handler,
}
