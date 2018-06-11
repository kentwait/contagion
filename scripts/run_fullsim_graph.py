#!/usr/bin/env python3
#
# run_fullsim_graph.py
#
# Generate directories for fixation_graph_fullsim and runs simulation
# In this experiment, the simulation runs on a graph
# under the SIS model with mutation and selection
# Variables:
# - network: reg, gnp, pow
# - mutation rate at 1e-4, 1e-5, 1e-6
# - Ns at 0, 1, 2, 4, 10
import argparse
import os
import pickle
import subprocess as proc
import numpy as np
import pandas as pd
from collections import Counter

DEFAULT_PATH = '/home/kent/data'
DEFAULT_BASE_RELPATH = 'fullsim_sis'
CONTAGION_PATH = 'contagion'
HOST_COUNT = 20

NETWORK_BASENAME = '{network}'
MU_BASENAME = '{network}_mu{mu:.9f}'
NS_BASENAME = '{network}_mu{mu:.9f}_Ns{Ns:+.0f}'

FASTA_FILENAME = 'pathogens.fa'
NETWORK_FILENAME = '{network}.list.fa'
FM_FILENAME = 'Ns{Ns:+.0f}.fm.txt'
TOML_FILENAME = 'config.{network}.mu{mu:.9f}.Ns{Ns:+.0f}.toml'
LOG_FILENAME = 'log.{network}.mu{mu:.9f}.Ns{Ns:+.0f}.txt'
DF_FILENAME = 'df.{network}.mu{mu:.9f}.Ns{Ns:+.0f}.factordf.pickle'
ARCHIVE_FILENAME = 'run.{network}.mu{mu:.9f}.Ns{Ns:+.0f}.tar.gz'
FREQ_CSV_FILENAME = 'log.001.freq.csv'
FREQ_DF_FILENAME = 'log.001.freq.df.pickle'
GENOTYPE_CSV_FILENAME = 'log.001.g.csv'
GENOTYPE_DF_FILENAME = 'log.001.g.df.pickle'

SEQUENCE = """111011010111101000100110010110110100010011110000100110111011
011001010111100111011110110110111100001000100111101111100101
011010111110010110010101110001101011111000010110010101111010
110001111110101000100000110101101000110101111111100010111101
101000111011110010100101100101100001110100001100110010000010
101101101011011001000100101111111001000101111010000001001001
001111010100110101010000110001001101111111011001101111000001
101111110010001000100001110111101000000000010000100000000011
110101010010111111100000011001011101111100011011110000110100
111010101111001000110010001101111010111101110000001101000100
010011010000001100001100010011001101111101110100000000110110
011000010010110110100011100010000001010001100011000000111101
011111101101011000110110101100001110111100011101000101110000
110110000010100010100100000000111100000010001010100101111100
001110010110101010000010100011100010010111110111000000100101
101111101001101100101010010100000110101001001110110101111000
1011011011001101011110111000101001110011"""

def generate_graph(path, graph_type):
    reg_g = """6	9	1.0
9	6	1.0
6	15	1.0
15	6	1.0
6	11	1.0
11	6	1.0
6	19	1.0
19	6	1.0
6	4	1.0
4	6	1.0
9	2	1.0
2	9	1.0
9	17	1.0
17	9	1.0
9	7	1.0
7	9	1.0
9	12	1.0
12	9	1.0
0	14	1.0
14	0	1.0
0	16	1.0
16	0	1.0
0	11	1.0
11	0	1.0
0	4	1.0
4	0	1.0
0	2	1.0
2	0	1.0
14	19	1.0
19	14	1.0
14	15	1.0
15	14	1.0
14	5	1.0
5	14	1.0
14	10	1.0
10	14	1.0
2	8	1.0
8	2	1.0
2	5	1.0
5	2	1.0
2	16	1.0
16	2	1.0
8	16	1.0
16	8	1.0
8	15	1.0
15	8	1.0
8	7	1.0
7	8	1.0
8	1	1.0
1	8	1.0
1	17	1.0
17	1	1.0
1	10	1.0
10	1	1.0
1	4	1.0
4	1	1.0
1	7	1.0
7	1	1.0
17	13	1.0
13	17	1.0
17	18	1.0
18	17	1.0
17	11	1.0
11	17	1.0
13	12	1.0
12	13	1.0
13	4	1.0
4	13	1.0
13	11	1.0
11	13	1.0
13	15	1.0
15	13	1.0
16	3	1.0
3	16	1.0
16	5	1.0
5	16	1.0
18	19	1.0
19	18	1.0
18	11	1.0
11	18	1.0
18	12	1.0
12	18	1.0
18	3	1.0
3	18	1.0
19	15	1.0
15	19	1.0
19	12	1.0
12	19	1.0
3	7	1.0
7	3	1.0
3	5	1.0
5	3	1.0
3	4	1.0
4	3	1.0
7	10	1.0
10	7	1.0
5	10	1.0
10	5	1.0
12	10	1.0
10	12	1.0"""
    gnp_g = """0	2	1.0
2	0	1.0
0	4	1.0
4	0	1.0
0	7	1.0
7	0	1.0
0	8	1.0
8	0	1.0
0	10	1.0
10	0	1.0
0	14	1.0
14	0	1.0
0	17	1.0
17	0	1.0
1	2	1.0
2	1	1.0
1	4	1.0
4	1	1.0
2	5	1.0
5	2	1.0
2	6	1.0
6	2	1.0
2	12	1.0
12	2	1.0
2	19	1.0
19	2	1.0
3	6	1.0
6	3	1.0
3	10	1.0
10	3	1.0
3	12	1.0
12	3	1.0
4	5	1.0
5	4	1.0
4	7	1.0
7	4	1.0
4	13	1.0
13	4	1.0
4	14	1.0
14	4	1.0
4	15	1.0
15	4	1.0
4	16	1.0
16	4	1.0
4	17	1.0
17	4	1.0
5	6	1.0
6	5	1.0
5	8	1.0
8	5	1.0
5	10	1.0
10	5	1.0
5	13	1.0
13	5	1.0
5	16	1.0
16	5	1.0
5	18	1.0
18	5	1.0
5	19	1.0
19	5	1.0
6	12	1.0
12	6	1.0
6	14	1.0
14	6	1.0
6	16	1.0
16	6	1.0
7	11	1.0
11	7	1.0
7	14	1.0
14	7	1.0
7	18	1.0
18	7	1.0
8	12	1.0
12	8	1.0
9	11	1.0
11	9	1.0
9	15	1.0
15	9	1.0
9	19	1.0
19	9	1.0
10	19	1.0
19	10	1.0
11	13	1.0
13	11	1.0
11	17	1.0
17	11	1.0
12	16	1.0
16	12	1.0
12	18	1.0
18	12	1.0
13	15	1.0
15	13	1.0
13	17	1.0
17	13	1.0
14	18	1.0
18	14	1.0
15	16	1.0
16	15	1.0
15	17	1.0
17	15	1.0
16	17	1.0
17	16	1.0"""
    pow_g = """0	3	1.0
3	0	1.0
0	4	1.0
4	0	1.0
0	5	1.0
5	0	1.0
0	7	1.0
7	0	1.0
0	11	1.0
11	0	1.0
0	13	1.0
13	0	1.0
0	18	1.0
18	0	1.0
1	3	1.0
3	1	1.0
1	4	1.0
4	1	1.0
1	5	1.0
5	1	1.0
1	8	1.0
8	1	1.0
1	9	1.0
9	1	1.0
1	10	1.0
10	1	1.0
1	15	1.0
15	1	1.0
1	19	1.0
19	1	1.0
2	3	1.0
3	2	1.0
2	6	1.0
6	2	1.0
2	9	1.0
9	2	1.0
2	14	1.0
14	2	1.0
3	4	1.0
4	3	1.0
3	5	1.0
5	3	1.0
3	6	1.0
6	3	1.0
3	8	1.0
8	3	1.0
3	9	1.0
9	3	1.0
3	10	1.0
10	3	1.0
3	12	1.0
12	3	1.0
3	14	1.0
14	3	1.0
3	15	1.0
15	3	1.0
3	16	1.0
16	3	1.0
3	17	1.0
17	3	1.0
4	7	1.0
7	4	1.0
4	11	1.0
11	4	1.0
4	18	1.0
18	4	1.0
5	7	1.0
7	5	1.0
5	8	1.0
8	5	1.0
5	11	1.0
11	5	1.0
5	12	1.0
12	5	1.0
5	13	1.0
13	5	1.0
7	13	1.0
13	7	1.0
8	12	1.0
12	8	1.0
8	15	1.0
15	8	1.0
8	16	1.0
16	8	1.0
8	17	1.0
17	8	1.0
9	10	1.0
10	9	1.0
9	14	1.0
14	9	1.0
9	19	1.0
19	9	1.0
10	18	1.0
18	10	1.0
12	16	1.0
16	12	1.0
14	19	1.0
19	14	1.0
16	17	1.0
17	16	1.0"""
    with open(path, 'w') as f:
        if graph_type == 'reg':
            print(reg_g, file=f)
        elif graph_type == 'gnp':
            print(gnp_g, file=f)
        elif graph_type == 'pow':
            print(pow_g, file=f)

def generate_sequence_fasta_multihost(path, infected_hostlist=[], npathogens=500):
    """Generates a FASTA file

    Parameters
    ----------
    path : str
        Location to save the sequences
    infected_hostlist : list
        IDs of infected hosts
    npathogens : int
        Number of sequences to generate
    """
    
    fasta_header = """# Equilibrium frequency
# No standing variation
% 0:0 1:1"""
    entry_format = '>h:{}\n{}'
    with open(path, 'w') as f:
        print(fasta_header, file=f)
        for i in infected_hostlist:
            for _ in range(npathogens):
                print(entry_format.format(i, SEQUENCE), file=f)

def generate_fm(path, Ns, N=500, nsites=1000):
    """Generates a multiplicative fitness matrix

    Parameters
    ----------
    path : str
        Location to save the fitness matrix
    Ns : float
        scaled selection coefficient
    N : int
        population size
    nsites : int
        Number of sites (length of sequence)
    """
    fm_header = """# Multiplicative fitness matrix
# This test models 1 site (index 0).
# This test uses 2 alleles per site.
# Sites ordered 0, 1
# Ns = {:+.0f}
# N = {}""".format(Ns, N)
    default_fm = 'default->{:+.8f}, 0.0'
    fm_format = '{}: {:+.8f}, 0.0'
    logF = np.log(1 - (Ns/float(N)))
    with open(path, 'w') as f:
        print(fm_header, file=f)
        print(default_fm.format(logF), file=f)
        print(fm_format.format(0, logF), file=f)
        if nsites > 1:
            print(fm_format.format(nsites-1, logF), file=f)

def generate_toml(path, pathogen_path, network_path, fm_path, csv_path, mu,
                  num_generations=10000, duration=10, transmission_prob=1.0, transmission_size=5, coinfection=False,
                  instances=1):
    template = """[simulation]
num_generations = {num_generations}
num_instances = {instances}
host_popsize = 20
epidemic_model = "sis"
coinfection = {coinfection}
pathogen_path = "{pathogen_path}"
host_network_path = "{network_path}"
num_sites = 1
expected_characters = ["0", "1"]

[logging]
log_freq = 1
log_path = "{csv_path}"
log_transmission = true

[[intrahost_model]]
model_name = "constant-nomutation"
host_ids = [
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19
]
mutation_rate = {mu}
transition_matrix = [
  [ 0.0e+00, 1.0e+00 ],
  [ 1.0e+00, 0.0e+00 ],
]
recombination_rate = 0.0
replication_model = "constant"
constant_pop_size = 500
infected_duration = {duration}

[[fitness_model]]
model_name = "multiplicative"
host_ids = [
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19
]
fitness_model = "multiplicative"
fitness_model_path = "{fm_path}"

[[transmission_model]]
model_name = "poisson"
host_ids = [
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19
]
mode = "poisson"
# each pathogen has a 0.01 chance of transmitting
transmission_prob = {transmission_prob:.1f}
transmission_size = {transmission_size:.1f}
"""
    with open(path, 'w') as f:
        print(template.format(
            mu=mu,
            pathogen_path=pathogen_path,
            network_path=network_path,
            fm_path=fm_path,
            csv_path=csv_path,
            num_generations=num_generations,
            duration=duration,
            transmission_prob=transmission_prob,
            transmission_size=transmission_size,
            coinfection='true' if coinfection else 'false',
            instances=instances,
        ), file=f)

def pickle_csv(csv_path, df_path):
    assert os.path.exists(csv_path) and os.path.isfile(csv_path), '%s is not a valid file path' % csv_path
    df = pd.read_csv(csv_path)
    pickle.dump(df, open(df_path, 'wb'))

def create_df(freq_path, genotype_path, instance, 
              samples=100, generation=1000,
              network=None, mu=None,Ns=None):
    # Read each realization and store value
    anc_seq = list(map(int, list(SEQUENCE.replace('\n', ''))))

    genotypes = pickle.load(open(genotype_path, 'rb'))
    freq = pickle.load(open(freq_path, 'rb'))
    freq_at_generation = freq[freq['generation'] == generation].sort_values(['hostID', 'freq'])
    
    # sampled_genotypes = []
    sampled_sequences = []
    host_ids = freq_at_generation['hostID'].unique()
    for host_id in host_ids[np.random.randint(len(host_ids), size=samples)]:
        df = freq_at_generation[freq_at_generation['hostID'] == host_id].sort_values(['freq']).reset_index()
        idx = np.argmax(np.random.multinomial(1, df.freq/df.freq.sum()))
        genotype_id = df.iloc[idx]['genotypeID']
        seq = list(map(int, str(genotypes[genotypes['genotypeID'] == genotype_id]['sequence'].values[0]).split()))
        # sampled_genotypes.append(genotype_id)
        sampled_sequences.append(seq)
    msa = np.array(sampled_sequences)
    
    # unfolded SFS
    diff = np.array(msa != np.array(anc_seq), dtype=np.int)  # Get number of derived alleles at each position
    counts = Counter(diff.T.sum(axis=1))
    num_sites = sum(counts.values())
    sfs_count = [{'class':k, 
                 'count':v, 
                 'instance':instance, 
                 'network':network,
                 'mutation_rate':mu,
                 'Ns':Ns,
                 } for k, v in counts.items()]
    return pd.DataFrame(sfs_count)
    
def run_contagion(toml_path, contagion_path=CONTAGION_PATH, threads=-1):
    cmd = [contagion_path, '-threads', str(threads), toml_path]
    p = proc.Popen(cmd, stdout=proc.PIPE, stderr=proc.STDOUT)
    while(True):
      retcode = p.poll() #returns None while subprocess is running
      line = p.stdout.readline()
      yield line
      if(retcode is not None):
        break

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('data_path', help='path to directory where data will be saved', type=str, default=DEFAULT_PATH)
    parser.add_argument('--base_relpath', help='basepath of the data relative to the data path', type=str, 
                        default=DEFAULT_BASE_RELPATH)
    parser.add_argument('--network', nargs='+', help='network type/s', type=str)
    parser.add_argument('--mu', nargs='+', help='mutation rate/s', type=float)
    parser.add_argument('--Ns', nargs='+', help='scaled selection coefficient/s Ns', type=float)
    parser.add_argument('--duration', help='duration of infection, default: 10', type=int, default=10)
    parser.add_argument('--transmission_size', help='transmission size, default: 5', type=int, default=5)
    parser.add_argument("--instances", help="number of realizations, default: 1000", type=int, default=1000)
    parser.add_argument("--generations", help="number of realizations, default: 1000", type=int, default=1000)
    parser.add_argument("--npathogens", help="number of pathogens, default: 500", type=int, default=500)
    parser.add_argument("--infected_nhosts", help="number of infected hosts in 20-host network, default: 1", type=int, default=1)
    parser.add_argument('--transmission_prob', help='transmission probability, default: 0.5', type=float, default=0.5)
    parser.add_argument("--coinfection", help="Allows coinfection", action='store_true')
    parser.add_argument("--sample_from_generation", help="generation number to sample from, default: 1000", type=int, default=1000)
    parser.add_argument("--sample_count", help="number to samples to pick for SFS, default: 100", type=int, default=100)
    parser.add_argument("--threads", help="number of threads to run Contagion, default: 2", type=int, default=2)
    parser.add_argument("--overwrite", help="overwrites existing files", action='store_true')
    parser.add_argument("--reduce_output", help="Only summaries are shown", action='store_true')
    parser.add_argument("--no_compression", help="Do not compress results", action='store_true')

    args = parser.parse_args()

    # Check if basepath exists
    basepath = os.path.abspath(os.path.join(args.data_path, args.base_relpath))
    assert os.path.exists(basepath), "Data base path does not exist. Check data_path ({}) and base_relpath ({}) values".format(args.data_path, args.base_relpath)

    # Assert that inputs are present
    assert args.network, 'No networks set (reg, gnp, pow)'
    assert args.duration, 'No durations set'
    assert args.transmission_size, 'No transmission sizes set'
    assert args.Ns, 'No Ns values set'    

    # data/fixation_graph_fullsim/graph
    for network in args.network:
        # Create directory
        network_dirpath = os.path.join(
            basepath, 
            NETWORK_BASENAME.format(network=network)
        )
        if not os.path.exists(network_dirpath):
            os.makedirs(network_dirpath)
        # Write network
        network_path = os.path.join(
            network_dirpath, 
            NETWORK_FILENAME.format(network=network)
        )
        generate_graph(network_path, network)

        # data/fixation_graph_bottleneck/graph/graph_mu0.000100000
        for mu in args.mu:
            # Create directory
            mu_path = os.path.join(
                network_dirpath, 
                MU_BASENAME.format(
                    network=network,
                    mu=mu,
                )
            )
            if not os.path.exists(mu_path):
                os.makedirs(mu_path)

            # data/fixation_graph_bottleneck/graph/graph_t005/graph_t005_m001/graph_t00d_m001_Ns+0
            for Ns in args.Ns:
                # Create directory
                ns_path = os.path.join(mu_path, NS_BASENAME.format(
                    network=network,
                    mu=mu, 
                    Ns=Ns,
                ))
                if os.path.exists(ns_path) and args.overwrite:
                    proc.call(['rm', '-R', ns_path])
                os.makedirs(ns_path)

                # Write fitness matrix file inside ns_path
                fm_path = os.path.join(ns_path, FM_FILENAME.format(Ns=Ns))
                generate_fm(fm_path, Ns, N=500, nsites=1)

                # Write log
                log_path = os.path.join(
                    ns_path, 
                    LOG_FILENAME.format(
                        network=network,
                        mu=mu, 
                        Ns=Ns,
                    )
                )
                # Overwrite previous log
                with open(log_path, 'w') as f:
                    print('', file=f)
                
                # create df list
                df_list = []

                # Create instances
                for instance in range(args.instances):
                    instance_path = os.path.join(
                        ns_path, 
                        '{:0>4d}'.format(instance)
                    )
                    os.makedirs(instance_path)

                    # Write sequences file
                    pathogen_path = os.path.join(instance_path, FASTA_FILENAME)
                    # Permute every instance
                    infected_hostlist = np.random.permutation(HOST_COUNT)[:args.infected_nhosts]
                    generate_sequence_fasta_multihost(
                        pathogen_path, 
                        infected_hostlist=infected_hostlist,
                        npathogens=args.npathogens, 
                    )

                    # Write TOML config file inside each instance_path
                    toml_path = os.path.join(
                        instance_path, 
                        TOML_FILENAME.format(
                            network=network,
                            mu=mu, 
                            Ns=Ns,
                        )
                    )
                    generate_toml(
                        toml_path, 
                        pathogen_path, 
                        network_path, 
                        fm_path, 
                        instance_path + '/',
                        mu,
                        num_generations=args.generations, 
                        duration=args.duration, 
                        transmission_prob=args.transmission_prob, 
                        transmission_size=args.transmission_size,
                        instances=1,
                        coinfection=args.coinfection,
                    )
                    # Run simulation
                    if not args.reduce_output:
                        print('Network:            {}'.format(network))
                        print('Mutation rate:      {}'.format(mu))
                        print('Ns:                 {:+.0f}'.format(Ns))
                        print('Instance:           {}'.format(instance))
                    with open(log_path, 'a') as f:
                        for output in run_contagion(toml_path, threads=args.threads, contagion_path=CONTAGION_PATH):
                            if not args.reduce_output:
                                print(output.decode('utf-8'), end='')
                            print(output.decode('utf-8'), end='', file=f)

                    # Pickle freq and genotype csv
                    # log.001.freq.csv
                    # log.001.g.csv
                    freq_path = os.path.join(instance_path, FREQ_CSV_FILENAME)
                    freq_pickle_path = os.path.join(
                        instance_path, FREQ_DF_FILENAME)
                    pickle_csv(freq_path, freq_pickle_path)

                    genotype_path = os.path.join(
                        instance_path, GENOTYPE_CSV_FILENAME)
                    genotype_pickle_path = os.path.join(
                        instance_path, GENOTYPE_DF_FILENAME)
                    pickle_csv(genotype_path, genotype_pickle_path)

                    # Create dataframe from freq and genotype df's
                    df = create_df(freq_pickle_path, genotype_pickle_path,  
                                   instance,
                                   network=network,
                                   mu=mu,
                                   Ns=Ns,
                                   samples=args.sample_count, 
                                   generation=args.sample_from_generation,
                    )
                    df_list.append(df)

                # Sample from freq and create pandas dataframe
                # Save df in ns_path
                df_pickle_path = os.path.join(
                    ns_path,
                    DF_FILENAME.format(
                        network=network,
                        mu=mu, 
                        Ns=Ns,
                    )
                )
                ns_df = pd.concat(df_list)
                pickle.dump(ns_df, open(df_pickle_path, 'wb'))
                # Save df in basepath
                df_pickle_path = os.path.join(
                    basepath,
                    DF_FILENAME.format(
                        network=network,
                        mu=mu, 
                        Ns=Ns,
                    )
                )
                ns_df = pd.concat(df_list)
                pickle.dump(ns_df, open(df_pickle_path, 'wb'))
                df_list = []

                if not args.no_compression:
                    # Archive folder
                    archive_path = os.path.join(
                        basepath, 
                        ARCHIVE_FILENAME.format(
                            network=network,
                            mu=mu, 
                            Ns=Ns,
                        )
                    )
                    proc.call(['tar', '-czf', archive_path, '-C', ns_path, '.'])
                    # Delete folder
                    proc.call(['rm', '-R', ns_path])

    print('Done.')
