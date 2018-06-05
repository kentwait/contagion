from collections import OrderedDict
import re
from prompt_toolkit import prompt
from prompt_toolkit.contrib.completers import WordCompleter


class Configuration(object):
    def __init__(self, config_path=None, contagion_path='contagion'):
        # Path to config file
        self.config_path = ''

        # simulation parameters
        self.num_generations = 0
        self.num_instances = 0
        self.host_popsize = 0
        self.epidemic_model = ''  # si, sir, sirs, sei, seis, seirs, endtrans, exchange
        self.coinfection = False
        # Path to pathogen sequences
        self.pathogen_path = ''
        # Path to adjacency list file (host network)
        self.host_network_path = ''

        # logging parameters
        self.log_freq = 0
        # Path to log
        self.log_path = ''

        # intrahost_model
        self.intrahost_model_dict = OrderedDict()

        # fitness_model
        self.fitness_model_dict = OrderedDict()

        # transmission_model
        self.transmission_model_dict = OrderedDict()

    def toml_string(self):
        config_string = ''

        # simulation
        config_string += '[simulation]\n'
        _param_dict = OrderedDict([
            ('num_generations', self.num_generations),
            ('num_instances', self.num_instances),
            ('host_popsize', self.host_popsize),
            ('epidemic_model', self.epidemic_model),
            ('coinfection', self.coinfection),
            ('pathogen_path', self.pathogen_path),
            ('host_network_path', self.host_network_path),
        ])
        for k,v in _param_dict.items():
            config_string += "{k} = {v}\n".format(k=k, v=v)
        config_string += '\n'

        # logging
        config_string += '[simulation]\n'
        _param_dict = OrderedDict([
            ('log_freq', self.log_freq),
            ('log_path', self.log_path),
        ])
        for k,v in _param_dict.items():
            config_string += "{k} = {v}\n".format(k=k, v=v)
        config_string += '\n'

        # intrahost model
        for model in self.intrahost_model_dict.values():
            config_string += model.toml_string()

        # fitness model
        for model in self.fitness_model_dict.values():
            config_string += model.toml_string()

        # transmission model
        for model in self.transmission_model_dict.values():
            config_string += model.toml_string()

        return config_string

    def save(self):
        with open(self.config_path, 'w') as f:
            print(self.toml_string(), file=f)

    def new_intrahost_model_prompt(self, history=None):
        """Creates a new intrahost_model object using answers from prompts.

        Parameters
        ----------
        history : prompt_toolkit.history.InMemoryHistory

        Returns
        -------
        configurator.configuration.IntrahostModel

        """
        # Initialize model
        model = IntrahostModel()

        # Add values
        # TODO: create name validator
        model.model_name = prompt(
            'Model name: ',
            history=history,
            validator=None,
        )
        host_ids = host_ids = prompt(
            'Host IDs: ',
            history=history,
            validator=None,
        )
        model.host_ids = parse_host_ids(host_ids)
        # Parameters
        model.mutation_rate = float(
            prompt('Mutation rate (subs/site/generation): ',
                   validator=None)
        )
        transition_matrix = prompt(
            'Conditioned transition rate matrix: ',
            validator=None,
        )
        model.transition_matrix = eval(transition_matrix)            
        model.transition_matrix = parse_transition_matrix(transition_matrix)
        model.recombination_rate = float(
            prompt('Recombination rate (recombinations/generation): ',
                   validator=None)
        )
        replication_model = prompt(
            'Replication model: ',
            validator=None,
            completer=WordCompleter(['constant', 'bht', 'fitness']),
        )
        # Model-dependent parameters
        if str(replication_model).lower() == 'constant':
            model.constant_pop_size = int(
                prompt('Population size: ',
                       validator=None)
            )
        elif str(replication_model).lower() == 'bht':
            model.max_pop_size = int(
                prompt('Maximum population size: ',
                       validator=None)
            )
            model.growth_rate = float(
                prompt('Growth rate: ',
                       validator=None)
            )
        elif str(replication_model).lower() == 'fitness':
            model.max_pop_size = int(
                prompt('Maximum population size: ',
                       validator=None)
            )
        # Durations
        if replication_model != 'fitness':
            # Check if epidemic model has been set
            if not self.epidemic_model:
                # Set epidemic model
                pass
            if self.epidemic_model == 'endtrans':
                model.infected_duration = int(
                    prompt('Duration at infected status: ',
                           validator=None)
                )
                model.removed_duration = int(
                    prompt('Duration at removed status: ',
                           validator=None)
                )
            elif self.epidemic_model == 'exchange':
                pass
            else:
                if self.epidemic_model.startswith('se'):
                    model.exposed_duration = int(
                        prompt('Duration at exposed status: ',
                               validator=None)
                    )
                if self.epidemic_model.startswith('si'):
                    model.infected_duration = int(
                        prompt('Duration at infected status: ',
                               validator=None)
                    )
                if self.epidemic_model.startswith('sei'):
                    model.infective_duration = int(
                        prompt('Duration at infective status: ',
                               validator=None)
                    )
                if self.epidemic_model.endswith('r'):
                    model.removed_duration = int(
                        prompt('Duration at removed status: ',
                               validator=None)
                    )
                if self.epidemic_model.endswith('rs'):
                    model.recovered_duration = int(
                        prompt('Duration at recovered status: ',
                               validator=None)
                    )
                # vaccinated_duration = int(
                #     prompt('Duration at vaccinated status: ',
                #            validator=None)
                # )
        return model

    def new_fitness_model_prompt(self, history=None):
        """Creates a new fitness_model object using answers from prompts.

        Parameters
        ----------
        history : prompt_toolkit.history.InMemoryHistory

        Returns
        -------
        configurator.configuration.FitnessModel

        """
        # Initialize model
        model = FitnessModel()

        # Add values
        model.model_name = prompt(
            'Model name: ',
            history=history,
            validator=None,
        )
        host_ids = host_ids = prompt(
            'Host IDs: ',
            history=history,
            validator=None,
        )
        model.host_ids = parse_host_ids(host_ids)
        model.fitness_model = prompt(
            'Fitness model: ',
            validator=None,
            completer=WordCompleter(['multiplicative', 'additive']),
        )
        # Ask to generate new fitness matrix or
        # pass a path to an existing matrix
        if str(prompt('Do you want to generate a fitness model? Y/n :',
                      default='Y')).lower() == 'y':
            model.new_fitness_matrix_prompt(
                model.fitness_model, 
                history=history,
            )
        else:
            model.fitness_model_path = prompt(
                'Fitness model path: ',
                history=history,
                validator=None,
            )
        return model

class Model(object):
    def __init__(self):
        self.model_name = ''
        self.host_ids = []


class IntrahostModel(Model):
    def __init__(self):
        super().__init__()
        self.mutation_rate = 0
        self.transition_matrix = []
        self.recombination_rate = 0
        self.replication_model = ''  # constant, bht, fitness
        # dependent on replication_model 
        self.constant_pop_size = 0
        self.max_pop_size = 0
        self.growth_rate = 0
        # Durations
        self.exposed_duration = 0
        self.infected_duration = 0
        self.infective_duration = 0
        self.removed_duration = 0
        self.recovered_duration = 0
        self.dead_duration = 0
        self.vaccinated_duration = 0

    def toml_string(self):
        config_string = '[[intrahost_model]]\n'
        _param_dict = OrderedDict([
            ('model_name', self.model_name),
            ('host_ids', self.host_ids),
            ('mutation_rate', self.mutation_rate),
            ('recombination_rate', self.recombination_rate),
            ('replication_model', self.replication_model),
        ])
        for k,v in _param_dict.items():
            config_string += "{k} = {v}\n".format(k=k, v=v)
        # Write durations only if value is non-zero
        _param_dict = OrderedDict([
            ('exposed_duration', self.exposed_duration),
            ('infected_duration', self.infected_duration),
            ('infective_duration', self.infective_duration),
            ('removed_duration', self.removed_duration),
            ('recovered_duration', self.recovered_duration),
            ('dead_duration', self.dead_duration),
            ('vaccinated_duration', self.vaccinated_duration),
        ])
        for k,v in _param_dict.items():
            if v != 0:
                config_string += "{k} = {v}\n".format(k=k, v=v)
        # Write model dependent params
        if self.replication_model in ['bht', 'fitness']:
            config_string += "{k} = {v}\n".format(
                k='max_pop_size', 
                v=self.max_pop_size,
            )
        else:
            config_string += "{k} = {v}\n".format(
                k='constant_pop_size', 
                v=self.constant_pop_size,
            )
        if self.replication_model in ['bh', 'bht']:
            config_string += "{k} = {v}\n".format(
                k='growth_rate', 
                v=self.growth_rate,
            )
        config_string += '\n'   


class FitnessModel(Model):
    def __init__(self):
        super().__init__()
        self.fitness_model = ''  # multiplicative, additive, additive_motif
        self.fitness_model_path = ''

    def toml_string(self):
        config_string = '[[fitness_model]]\n'
        _param_dict = OrderedDict([
            ('model_name', self.model_name),
            ('host_ids', self.host_ids),
            ('fitness_model', self.fitness_model),
            ('fitness_model_path', self.fitness_model_path),
        ])
        for k,v in _param_dict.items():
            config_string += "{k} = {v}\n".format(k=k, v=v)
        config_string += '\n'

    def new_fitness_matrix_prompt(self, fitness_model, history=None):
        """Wizard that creates a new fitness matrix based on answers
        to prompts.

        Parameters
        ----------
        fitness_model : str
        history : prompt_toolkit.history.InMemoryHistory

        """
        self.fitness_model_path = prompt(
            'Fitness model save path: ',
            history=history,
            validator=None,
        )
        self.fitness_model = fitness_model
        # generate fitness matrix
        if str(prompt('Create neutral matrix [Y/n]: ',
                      default='Y')).lower() == 'y':
            num_sites = int(prompt('Number of sites: ', validator=None))
            num_variants = int(
                prompt('Number of potential states per site: ',
                       validator=None),
            )
            if self.fitness_model == 'multiplicative':
                self.generate_neutral_matrix(
                    num_sites,
                    num_variants,
                    self.fitness_model_path,
                )
            else:
                growth_rate = float(prompt('Growth rate: ', validator=None))
                self.generate_additive_neutral_matrix(
                    num_sites,
                    num_variants,
                    growth_rate,
                    self.fitness_model_path,
                )
        else:
            num_sites = int(prompt('Number of sites: ', validator=None))
            fitnesses = prompt('Enter list of fitness values: ', validator=None)
            fitnesses = eval(fitnesses)
            if self.fitness_model == 'multiplicative':
                self.generate_single_preference_matrix(
                    num_sites,
                    fitnesses,
                    self.fitness_model_path,
                )
            else:
                growth_rates = prompt('Enter list of growth rates: ')
                self.generate_additive_single_preference_matrix(
                    num_sites,
                    growth_rates,
                    self.fitness_model_path,
                )

    @staticmethod
    def generate_neutral_matrix(num_sites, num_variants, save_path):
        """Generates a multiplicative neutral fitness matrix 
        and writes it to file.

        Parameters
        ----------
        num_sites : int
        num_variants : int
        save_path : str

        """
        fitness_values = ', '.join(['1.0' for _ in range(num_variants)])
        text = 'default->' + fitness_values + '\n'
        text += '0: ' + fitness_values + '\n'
        text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
        with open(save_path, 'w') as f:
            print(text, file=f)

    @staticmethod
    def generate_additive_neutral_matrix(num_sites, num_variants, growth_rate, save_path):
        """Generates an additive neutral fitness matrix and writes it to file.

        Parameters
        ----------
        num_sites : int
        num_variants : int
        growth_rate : float
        save_path : str

        """
        fitness_values = ', '.join([str(growth_rate/num_sites) for _ in range(num_variants)])
        text = 'default->' + fitness_values + '\n'
        text += '0: ' + fitness_values + '\n'
        text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
        with open(save_path, 'w') as f:
            print(text, file=f)

    @staticmethod
    def generate_single_preference_matrix(num_sites, fitnesses, save_path):
        """Generates a multiplicative single preference fitness matrix
        and writes it to file.

        Parameters
        ----------
        num_sites : int
        fitnesses : list of float
        save_path : str

        """
        fitness_values = ', '.join(
            [str(f) for f in map(float, re.findall(r'\d*\.?\d+', fitnesses))]
        )
        text = 'default->' + fitness_values + '\n'
        text += '0: ' + fitness_values + '\n'
        text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
        with open(save_path, 'w') as f:
            print(text, file=f)

    @staticmethod
    def generate_additive_single_preference_matrix(num_sites, growth_rates, save_path):
        """Generates an additive single preference fitness matrix
        and writes it to file.

        Parameters
        ----------
        num_sites : int
        growth_rates : list of float
        save_path : str

        """
        fitness_values = ', '.join(
            [str(f/num_sites) for f in map(float, re.findall(r'\d*\.?\d+', growth_rates))]
        )
        text = 'default->' + fitness_values + '\n'
        text += '0: ' + fitness_values + '\n'
        text += '{}: '.format(num_sites - 1) + fitness_values + '\n'
        with open(save_path, 'w') as f:
            print(text, file=f)


class TransmissionModel(Model):
    def __init__(self):
        super().__init__()
        self.mode = ''  # poisson, constant
        self.transmission_prob = 0
        self.transmission_size = 0

    def toml_string(self):
        config_string = '[[transmission_model]]\n'
        _param_dict = OrderedDict([
            ('model_name', self.model_name),
            ('host_ids', self.host_ids),
            ('transmission_prob', self.transmission_prob),
            ('transmission_size', self.transmission_size),
        ])
        for k,v in _param_dict.items():
            config_string += "{k} = {v}\n".format(k=k, v=v)
        config_string += '\n' 


def parse_host_ids(text):
    """Parses the host list input string. If preceded by a "!", generate a host list
    using the input as a range.

    Parameters
    ----------
    test : str

    Returns
    -------
    list

    """
    match = re.search(r'^\!\[\s*(\d+)\s*\,\s*(\d+)\s*\,?\s*(\d+)?\s*\,?\]$', text)
    if match:
        if len(match.groups()) == 3:
            start, end, skip = list(map(int, match.groups()))
        else:
            start, end = list(map(int, match.groups()))
            skip = 1
        return [i for i in range(start, end, skip)]
    return list(map(int, re.findall(r'\d+', text)))
