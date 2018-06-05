from configurator.configuration import Configuration
from configurator import handlers

EXIT_COMMANDS = ['exit', 'exit()', 'quit', 'quit()', 'q']
SINGLE_WORD_COMMANDS = EXIT_COMMANDS + ['configure', 'clear']
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
