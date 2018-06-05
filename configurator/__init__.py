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
