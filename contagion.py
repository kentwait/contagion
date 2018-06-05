from prompt_toolkit.history import InMemoryHistory
from prompt_toolkit.contrib.completers import WordCompleter
from prompt_toolkit import prompt
import networkx as nx
import sys
import os
import re

from configurator import Configuration
from configurator import handlers
from configurator.configuration import IntrahostModel, FitnessModel, TransmissionModel
from configurator.validators import DirExistsValidator, StatementValidator
from configurator.sequence_generator import multiple_host_random_fasta as multihost

PROMPT = 'contagion> '
EXIT_COMANNDS = ['exit', 'exit()', 'quit', 'quit()', ]
SINGLEVALUE_COMMANDS = EXIT_COMANNDS + ['configure',]
SPECIAL_COMMANDS = ['create', 'reset']
OP_COMMANDS = ['set', ]
COMMANDS = OP_COMMANDS + EXIT_COMANNDS
CONFIG_PROPERTIES = list(Configuration().__dict__.keys())

def prompt_config_path():  
    default_path = os.path.join(os.getcwd(), 'config.toml')
    config_path = prompt(
        'Configuration file path [{}]: '.format(default_path), 
        validator=DirExistsValidator()
    )
    if not config_path:
        config_path = default_path
    return config_path

def configuration_wizard(config_obj, history=None):
    
    num_generations = prompt('Number of generations: ')
    num_instances = prompt('Number of simulation trials: ')
    host_popsize = prompt('Host population size: ')
    host_popsize = int(host_popsize)
    epidemic_model_completer = WordCompleter([
        'si', 'sir', 'sirs', 
        'sei', 'seis', 'seirs', 
        'endtrans', 'exchange'
    ])
    epidemic_model = prompt('Epidemic model: ', completer=epidemic_model_completer)
    coinfection = prompt('Allow coinfection [y/N]: ', default='N')

    # Generate pathogens
    if prompt('Do you want to generate pathogen sequences? [Y/n]: ', default='Y'):
        pathogen_path = prompt('Pathogen sequences save path: ', history=history, validator=None)
        num_sites = prompt('Length of pathogen sequence: ')
        chars = prompt('Character states: ')
        popsize = prompt('Number of pathogens per host: ')
        char_dist = prompt('Probability of picking a character (ENTER if uniform probability): ')
        char_dist = [float(f) for f in re.findall(r'\d*\.?\d+', char_dist)]

        host_ids = prompt('Infected host IDs: ', history=history, validator=None)
        host_ids = parse_host_ids(host_ids)
        clonal_completer = WordCompleter([
            'all', 'host', 'random', 
        ])
        clonal = prompt('Pathogen sequence identity (all|host|random): ', completer=clonal_completer)
        # Generate
        fasta_text = multihost(host_ids, int(popsize), chars, 
        int(num_sites), clonal, char_dist)
        with open(pathogen_path, 'w') as f:
            print(fasta_text, file=f)

    else:
        pathogen_path = prompt('Pathogen sequences path: ', history=history, validator=None)

    # Generate network
    if prompt('Do you want to generate a random network? [Y/n]: ', default='Y'):
        host_network_path = prompt('Host network save path: ', history=history, validator=None)
        network_completer = WordCompleter([
            'gnp', 'binomial', 'erdos-renyi', 
            'barabasi-albert', 'scale-free', 
            'holme-kim', 'powerlaw_cluster',
            'complete',
        ])
        while True:
            network_type = prompt('Network type: ', completer=network_completer)
            n = host_popsize
            if str(network_type).lower() in ['gnp', 'binomial', 'erdos-renyi']:
                # gnp
                # n (int) – The number of nodes.
                # p (float) – Probability for edge creation.
                p = prompt('Probability of an edge between two nodes p: ')
                network = nx.fast_gnp_random_graph(n, float(p), directed=False)
                with open(host_network_path, 'w') as f:
                    for a, row in network.adjacency():
                        for b in row:
                            print(a, b, 1.0, file=f)
                break

            elif str(network_type).lower() in ['barabasi-albert', 'scale-free']:
                # ba
                # n (int) – Number of nodes
                # m (int) – Number of edges to attach from a new node to existing nodes
                # error if m does not satisfy 1 <= m < n
                m = prompt('Number of edges to attach from a new node to existing nodes m (1 <= m < n): ')
                network = nx.barabasi_albert_graph(n, int(m))
                with open(host_network_path, 'w') as f:
                    for a, row in network.adjacency():
                        for b in row:
                            print(a, b, 1.0, file=f)
                break
                
            elif str(network_type).lower() in ['holme-kim', 'powerlaw_cluster']:
                #hk
                # n (int) – the number of nodes
                # m (int) – the number of random edges to add for each new node
                # p (float,) – Probability of adding a triangle after adding a random edge
                # error if m does not satisfy 1 <= m <= n or p does not satisfy 0 <= p <= 1
                m = prompt('Number of random edges to add for each new node m (1 <= m <= n): ')
                p = prompt('Probability of adding a triangle after adding a random edge p (0 <= p <= 1): ')
                network = nx.powerlaw_cluster_graph(n, int(m), float(p))
                with open(host_network_path, 'w') as f:
                    for a, row in network.adjacency():
                        for b in row:
                            print(a, b, 1.0, file=f)
                break
                
            elif str(network_type).lower() == 'complete':
                # complete
                network = nx.complete_graph(n)
                with open(host_network_path, 'w') as f:
                    for a, row in network.adjacency():
                        for b in row:
                            print(a, b, 1.0, file=f)
                break
                
            else:
                print('unrecognized network type')

    else:
        host_network_path = prompt('Host network path: ', history=history, validator=None)

    config_obj.num_generations = int(num_generations)
    config_obj.num_instances = int(num_instances)
    config_obj.host_popsize = int(host_popsize)
    config_obj.epidemic_model = epidemic_model
    config_obj.coinfection = bool(coinfection)
    config_obj.pathogen_path = pathogen_path
    config_obj.host_network_path = host_network_path

    num_intrahost_model = prompt('How many intrahost models do you want to create: ')
    for i in range(int(num_intrahost_model)):
        create_intrahost_model(epidemic_model, history=history)
    num_fitness_model = prompt('How many fitness models do you want to create: ')
    for i in range(int(num_fitness_model)):
        create_fitness_model(history=history)
    num_transmission_model = prompt('How many transmission models do you want to create: ')
    for i in range(int(num_transmission_model)):
        create_transmission_model(history=history)

    if str(prompt('Do you want to save this configuration? [Y/n]: ', default='Y')).lower() == 'y':
        config_obj.save()

def create_model(config_obj, text, history=None):
    # Add model to config
    model = text.split()[1]
    if model == 'transmission_model':
        model = create_transmission_model(history)
        config_obj.transmission_model_dict[model.model_name] = model
    elif model == 'fitness_model':
        model = create_fitness_model(history)
        config_obj.fitness_model_dict[model.model_name] = model
    elif model == 'intrahost_model':
        epidemic_model = config_obj.epidemic_model
        if not config_obj.epidemic_model:
            epidemic_model = set_epidemic_model(config_obj, history=history)
        model = create_intrahost_model(epidemic_model, history)
        config_obj.fitness_model_dict[model.model_name] = model

def create_intrahost_model(epidemic_model, history=None):
    replication_model_completer = WordCompleter(['constant', 'bht', 'fitness'])

    model_name = prompt('Model name: ', history=history, validator=None)
    host_ids = prompt('Host IDs: ', history=history, validator=None)
    host_ids = parse_host_ids(host_ids)

    mutation_rate = prompt('Mutation rate (subs/site/generation): ', history=history, validator=None)
    # transmission matrix
    transition_matrix = prompt('Conditioned transition rate matrix: ', history=history, validator=None)

    recombination_rate = prompt('Recombination rate (recombinations/generation): ', history=history, validator=None)
    replication_model = prompt('Replication model: ', history=history, validator=None, completer=replication_model_completer)

    # model dependent params
    if str(replication_model).lower == 'constant':
        constant_pop_size = prompt('Population size: ', history=history, validator=None)
    if str(replication_model).lower == 'bht' or str(replication_model).lower == 'fitness':
        max_pop_size = prompt('Maximum population size: ', history=history, validator=None)
    if str(replication_model).lower == 'bht':
        growth_rate = prompt('Growth rate: ', history=history, validator=None)

    # durations
    exposed_duration = 0
    infected_duration = 0
    infective_duration = 0
    removed_duration = 0
    recovered_duration = 0
    if replication_model != 'fitness':
        if epidemic_model == 'endtrans':
            infected_duration = prompt('Duration at infected status: ', history=history, validator=None)
            removed_duration = prompt('Duration at removed status: ', history=history, validator=None)
        elif epidemic_model == 'exchange':
            pass
        else:
            if epidemic_model.startswith('se'):
                exposed_duration = prompt('Duration at exposed status: ', history=history, validator=None)
            if epidemic_model.startswith('si'):
                infected_duration = prompt('Duration at infected status: ', history=history, validator=None)
            if epidemic_model.startswith('sei'):
                infective_duration = prompt('Duration at infective status: ', history=history, validator=None)
            if epidemic_model.endswith('r'): 
                removed_duration = prompt('Duration at removed status: ', history=history, validator=None)
            if epidemic_model.endswith('rs'):    
                recovered_duration = prompt('Duration at recovered status: ', history=history, validator=None)
            # vaccinated_duration = prompt('Duration at vaccinated status: ', history=history, validator=None)

    # Create model
    model = IntrahostModel()
    model.model_name = model_name
    model.host_ids = host_ids
    model.mutation_rate = float(mutation_rate)
    # TODO: replace eval with a function
    model.transition_matrix = eval(transition_matrix)
    model.recombination_rate = float(recombination_rate)
    model.replication_model = replication_model
    # model dependent params
    if str(replication_model).lower == 'constant':
        model.constant_pop_size = int(constant_pop_size)
    if str(replication_model).lower == 'bht' or str(replication_model).lower == 'fitness':
        model.max_pop_size = int(max_pop_size)
    if str(replication_model).lower == 'bht':
        model.growth_rate = float(growth_rate)
    # Durations
    model.exposed_duration = int(exposed_duration)
    model.infected_duration = int(infected_duration)
    model.infective_duration = int(infective_duration)
    model.removed_duration = int(removed_duration)
    model.recovered_duration = int(recovered_duration)
    # model.vaccinated_duration = vaccinated_duration

    return model

def create_fitness_model(history=None):
    fitness_model_completer = WordCompleter(['multiplicative', 'additive'])

    model_name = prompt('Model name: ', history=history, validator=None)
    host_ids = prompt('Host IDs: ', history=history, validator=None)
    host_ids = parse_host_ids(host_ids)
    fitness_model = prompt('Fitness model: ', history=history, validator=None, completer=fitness_model_completer)

    # Generate fitness model or pass existing
    generate_model = prompt('Do you want to generate a fitness model? Y/n :', default='Y')
    if str(generate_model).lower() == 'y':
        fitness_model_path = prompt('Fitness model save path: ', history=history, validator=None)
        # TODO: add validator for path
        # generate fitness model
        if str(prompt('Create neutral model [Y/n]: ', default='Y')).lower() == 'y':
            num_sites = prompt('Number of sites: ')
            num_variants = prompt('Number of potential states per site: ')
            if fitness_model == 'multiplicative':
                generate_neutral_fitness(int(num_sites), int(num_variants), fitness_model_path)
            else:
                growth_rate = prompt('Growth rate: ')
                generate_additive_neutral_fitness(int(num_sites), int(num_variants), float(growth_rate), fitness_model_path)
        else:
            num_sites = prompt('Number of sites: ')
            fitnesses = prompt('Enter list of fitness values: ')
            if fitness_model == 'multiplicative':
                generate_unipreference_fitness(int(num_sites), fitnesses, fitness_model_path)
            else:
                growth_rates = prompt('Enter list of growth rates: ')
                generate_additive_unipreference_fitness(int(num_sites), growth_rates, fitness_model_path)

    else:
        fitness_model_path = prompt('Fitness model path: ', history=history, validator=None)

    # Create model
    model = FitnessModel()
    model.model_name = model_name
    model.host_ids = host_ids
    model.fitness_model = fitness_model
    model.fitness_model_path = fitness_model_path
    return model

def create_transmission_model(history=None):
    model_name = prompt('Model name: ', history=history, validator=None)
    host_ids = prompt('Host IDs: ', history=history, validator=None)
    host_ids = parse_host_ids(host_ids)
    transmission_prob = prompt('Transmission probability: ', history=history, validator=None)
    transmission_size = prompt('Transmission size: ', history=history, validator=None)
    # Create model
    model = TransmissionModel()
    model.model_name = model_name
    model.host_ids = host_ids
    model.transmission_prob = transmission_prob
    model.transmission_size = transmission_size
    return model

def set_epidemic_model(config_obj, history=None):
    models = [
        'si', 'sis', 'sir', 'sirs', 'sei', 'seir', 'seirs',
        'endtrans', 'exchange',    
    ]
    model_completer = WordCompleter(models)
    text = prompt('Epidemic model: ', history=history, validator=None, completer=model_completer)
    if text:
        config_obj.epidemic_model = text
    return text

def set_property(config_obj, text):
    _, kv = text.split(None, 1)
    name, value = re.split(r'\s*\=\s*', kv)
    value = value.lstrip("'").lstrip('"').rstrip("'").rstrip('"')
    config_obj.__setattr__(name, value)

    # TODO: return a message confirming that the value was set

def return_property(config_obj, text):
    if text in config_obj.__dict__:
        v = config_obj.__getattribute__(text)
        if v is None:
            return 'None'
        elif isinstance(v, str) and v == '':
            return "''"
        return v
    return 'unknown configuration parameter'

def parse_host_ids(text):
    generate_match = re.search(r'^\!\[\s*(\d+)\s*\,\s*(\d+)\s*\,?\s*(\d+)\s*?\,?\]$', text)
    if generate_match:
        start, end, skip = generate_match.groups()
        return repr([i for i in range(int(start), int(end), int(skip))])
    return map(int, re.findall(r'\d+', text))

def generate_neutral_fitness(num_sites, num_variants, save_path):
    fitness_values = ', '.join(['1.0' for _ in range(num_variants)])
    text = 'default->' + fitness_values + '\n'
    text += '0: ' + fitness_values + '\n'
    text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
    with open(save_path, 'w') as f:
        print(text, file=f)

def generate_additive_neutral_fitness(num_sites, num_variants, growth_rate, save_path):
    fitness_values = ', '.join([str(growth_rate/num_sites) for _ in range(num_variants)])
    text = 'default->' + fitness_values + '\n'
    text += '0: ' + fitness_values + '\n'
    text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
    with open(save_path, 'w') as f:
        print(text, file=f)

def generate_unipreference_fitness(num_sites, fitnesses, save_path):
    fitness_values = ', '.join([str(f) for f in map(float, re.findall(r'\d*\.?\d+', fitnesses))])
    text = 'default->' + fitness_values + '\n'
    text += '0: ' + fitness_values + '\n'
    text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
    with open(save_path, 'w') as f:
        print(text, file=f)
        
def generate_additive_unipreference_fitness(num_sites, growth_rates, save_path):
    fitness_values = ', '.join([str(f/num_sites) for f in map(float, re.findall(r'\d*\.?\d+', growth_rates))])
    text = 'default->' + fitness_values + '\n'
    text += '0: ' + fitness_values + '\n'
    text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
    with open(save_path, 'w') as f:
        print(text, file=f)

EXIT_COMMANDS = ['exit', 'exit()', 'quit', 'quit()', 'q']
SINGLE_WORD_COMMANDS = EXIT_COMANNDS + ['configure', 'clear']
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

def main(config_path=None, contagion_path='contagion'):
    # Create configuration object
    config_obj = Configuration(config_path=config_path, contagion_path=contagion_path)
    # Instantiate history
    history = InMemoryHistory()
    # shell interface loop
    while True:
        try:
            # Valid statements:
            # configure
            # run logger=<csv|sqlite> threads=<int>
            # create intrahost_model|fitness_model|transmission_model <model_name>
            # append intrahost_model|fitness_model|transmission_model <model_name>
            # generate pathogens|network|fitness_matrix
            # set <configuration property>=<value>
            # get <configuration property>
            # reset <configuration property>
            # load configuration <path>
            # save configuration <path>
            # todb <basepath> <outpath>
            # tocsv <path> <basepath>
            # exit|exit()|quit|quit()|q
            # clear
            text = prompt(PROMPT, history=history, validator=StatementValidator())
        except KeyboardInterrupt:  # Ctrl+C
            continue
        except EOFError:  # Ctrl+D
            break
        else:
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
                elif text.split(None, 1) in PREFIX_COMMAND_HANDLER.keys():
                    args = [arg for arg in text.split(None) if '=' not in arg]
                    kwargs = [kwarg for kwarg in text.split(None) if '=' in kwarg]
                    PREFIX_COMMAND_HANDLER[text.split(None, 1)](*args, config_obj=config_obj, **kwargs)

if __name__ == '__main__':
    # TODO: Use click to get contagion path
    contagion_path = '/Volumes/Data/golang/src/github.com/kentwait/contagiongo/contagion'
    config_path = sys.argv[1] if len(sys.argv) > 1 else None
    main(config_path, contagion_path)