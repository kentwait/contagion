#!/usr/bin/env python3
#
# run_fixation_graph_bottleneck.py
#
# Generate directories for fixation_graph_bottleneck and runs simulation
# In this experiment, the simulation runs on a graph
# under the SIS model to determine the fixation probability
# Variables:
# - network: reg, gnp, pow
# - duration of infection: 5, 10, 20, 40 generations
# - transmission size: 1, 5, 50, 500 pathogens (constant) or 1, 1%, 10%, 100% bottleneck
# - Ns at 0, 1, 2, 4, 10
import numpy as np
import argparse
import os
import subprocess as proc
import re
import numpy as np

DEFAULT_PATH = '/home/kent/data'
DEFAULT_BASE_RELPATH = 'fixation_graph_bottleneck'
CONTAGION_PATH = 'contagion'
HOST_COUNT = 20

NETWORK_BASENAME = '{network}'
DURATION_BASENAME = '{network}_t{duration:0>2d}'
TSIZE_BASENAME = '{network}_t{duration:0>2d}_m{tsize:0>3d}'
NS_BASENAME = '{network}_t{duration:0>2d}_m{tsize:0>3d}_Ns{Ns:+.0f}'

FASTA_FILENAME = 'pathogens.fa'
NETWORK_FILENAME = '{network}.list.fa'
FM_FILENAME = 'Ns{Ns:+.0f}.fm.txt'
TOML_FILENAME = 'config.{network}.t{duration:0>2d}.m{tsize:0>3d}.Ns{Ns:+.0f}.toml'
LOG_FILENAME = 'log.{network}.t{duration:0>2d}.m{tsize:0>3d}.Ns{Ns:+.0f}.txt'
SUMMARY_FILENAME = 'summary.{network}.t{duration:0>2d}.m{tsize:0>3d}.Ns{Ns:+.0f}.txt'
ARCHIVE_FILENAME = 'run.{network}.t{duration:0>2d}.m{tsize:0>3d}.Ns{Ns:+.0f}.tar.gz'

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

def generate_single_site_fasta_multihost(path, infected_hostlist=[], npathogens=500, p=0.5):
    """Generates a FASTA file for single sites

    Parameters
    ----------
    path : str
        Location to save the sequences
    infected_hostlist : list
        IDs of infected hosts
    npathogens : int
        Number of sequences to generate
    p : float
        Initial frequency of allele. Probability of 0 or 1.
    """
    header = """# Single-site sequence
# Linear network configuration
% 0:0 1:1"""
    entry_format = '>h:{host}\n{seq}'
    with open(path, 'w') as f:
        print(header, file=f)
        for i in infected_hostlist:
            for _ in range(0, int(npathogens*p)):
                seq = 0
                print(entry_format.format(host=i, seq=seq), file=f)
            for _ in range(int(npathogens*p), npathogens):
                seq = 1
                print(entry_format.format(host=i, seq=seq), file=f)

def generate_fm(path, Ns, N=500, nsites=1):
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
            print(fm_format.format(nsites, logF), file=f)

def generate_toml(path, pathogen_path, network_path, fm_path, csv_path,
                  num_generations=10000, duration=10, transmission_prob=1.0, transmission_size=5,
                  instances=1):
    template = """[simulation]
num_generations = {num_generations}
num_instances = {instances}
host_popsize = 20
epidemic_model = "sis"
coinfection = false
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
mutation_rate = 0.0
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
model_name = "constant"
host_ids = [
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10, 11, 12, 13, 14, 15, 16, 17, 18, 19
]
mode = "poisson"
# each pathogen has a 0.01 chance of transmitting
transmission_prob = {transmission_prob:.1f}
transmission_size = {transmission_size:.1f}

[[stop_condition]]
condition = "allele_fixloss"
sequence = "1"
position = 0
"""
    with open(path, 'w') as f:
        print(template.format(
            pathogen_path=pathogen_path,
            network_path=network_path,
            fm_path=fm_path,
            csv_path=csv_path,
            num_generations=num_generations,
            duration=duration,
            transmission_prob=transmission_prob,
            transmission_size=transmission_size,
            instances=instances,
        ), file=f)

def summarize_fixation_graph(summary_path, log_path, Ns, network, reversed_values=False):
    assert os.path.exists(log_path), 'Log file in log_path does not exist'

    fixed = 0
    lost = 0
    iters = 0
    lost_flag = False
    fixed_flag = False
    lost_gens = []
    fixed_gens = []
    with open(log_path, 'r') as f:
        for line in f.readlines():
            if 'allele lost' in line:
                lost += 1
                lost_flag = True
                continue
            elif 'allele fixed' in line:
                fixed += 1
                fixed_flag = True
                continue
            elif 'Finished' in line:
                iters += 1

            if lost_flag:
                match = re.search(r'\d+', line)
                if match:
                    lost_gens.append(int(match.group()))
                lost_flag = False
            elif fixed_flag:
                match = re.search(r'\d+', line)
                if match:
                    fixed_gens.append(int(match.group()))
                fixed_flag = False
    
    with open(summary_path, 'w') as f:
        print('-'*80)
        if reversed_values:
            print(os.path.abspath(log_path))
            print('network: {}'.format(network))
            print('Ns:      {:+.0f}'.format(Ns))
            print('fixed:   {}\t{:.4f}'.format(lost, lost/float(iters)))
            if len(lost_gens) != 0:
                print('mean tf: {:.4f} generations'.format(np.mean(lost_gens)))
            else:
                print('mean tf: None')
            print('lost:    {}\t{:.4f}'.format(fixed, fixed/float(iters)))
            print('trials:  {}'.format(iters))

            print(os.path.abspath(log_path), file=f)
            print('network: {}'.format(network), file=f)
            print('Ns:      {:+.0f}'.format(Ns), file=f)
            print('fixed:   {}\t{:.4f}'.format(lost, lost/float(iters)), file=f)
            if len(lost_gens) != 0:
                print('mean tf: {:.4f} generations'.format(np.mean(lost_gens)), file=f)
            else:
                print('mean tf: None', file=f)
            print('lost:    {}\t{:.4f}'.format(fixed, fixed/float(iters)), file=f)
            print('trials:  {}'.format(iters), file=f)
        else:
            print(os.path.abspath(log_path))
            print('network: {}'.format(network))
            print('Ns:      {:+.0f}'.format(Ns))
            print('fixed:   {}\t{:.4f}'.format(fixed, fixed/float(iters)))
            if len(fixed_gens) != 0:
                print('mean tf: {:.4f} generations'.format(np.mean(fixed_gens)))
            else:
                print('mean tf: None', file=f)
            print('lost:    {}\t{:.4f}'.format(lost, lost/float(iters)))
            print('trials:  {}'.format(iters))

            print(os.path.abspath(log_path), file=f)
            print('network    {}'.format(network), file=f)
            print('Ns:      {:+.0f}'.format(Ns), file=f)
            print('fixed:   {}\t{:.4f}'.format(fixed, fixed/float(iters)), file=f)
            if len(fixed_gens) != 0:
                print('mean tf: {:.4f} generations'.format(np.mean(fixed_gens)), file=f)
            else:
                print('mean tf: None', file=f)
            print('lost:    {}\t{:.4f}'.format(lost, lost/float(iters)), file=f)
            print('trials:  {}'.format(iters), file=f)
        print('-'*80)

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
    parser.add_argument('--duration', nargs='+', help='duration/s of infection', type=int)
    parser.add_argument('--transmission_size', nargs='+', help='transmission size/s', type=int)
    parser.add_argument('--Ns', nargs='+', help='scaled selection coefficient/s Ns', type=float)
    parser.add_argument("--p", help="initial frequency", type=float, default=0.5)
    parser.add_argument("--instances", help="number of realizations", type=int, default=1000)
    parser.add_argument("--generations", help="number of realizations", type=int, default=10000)
    parser.add_argument("--npathogens", help="number of pathogens", type=int, default=500)
    parser.add_argument("--infected_nhosts", help="number of infected hosts in 20-host network", type=int, default=1)
    parser.add_argument('--transmission_prob', help='transmission probability, default: 1.0', type=float, default=1.0)
    parser.add_argument("--threads", help="number of threads to run Contagion", type=int, default=2)
    parser.add_argument("--overwrite", help="overwrites existing files", action='store_true')
    parser.add_argument("--reduce_output", help="Only summaries are shown", action='store_true')

    args = parser.parse_args()

    # Assert that inputs are present
    assert args.network, 'No networks set (reg, gnp, pow)'
    assert args.duration, 'No durations set'
    assert args.transmission_size, 'No transmission sizes set'
    assert args.Ns, 'No Ns values set'

    # Check if basepath exists
    basepath = os.path.abspath(os.path.join(args.data_path, args.base_relpath))
    assert os.path.exists(basepath), "Data base path does not exist. Check data_path ({}) and base_relpath ({}) values".format(args.path, args.base_relpath)

    # data/fixation_graph_bottleneck/graph
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

        # data/fixation_graph_bottleneck/graph/graph_t005
        for duration in args.duration:
            # Create directory
            duration_path = os.path.join(
                network_dirpath, 
                DURATION_BASENAME.format(
                    network=network,
                    duration=duration,
                )
            )
            if not os.path.exists(duration_path):
                os.makedirs(duration_path)

            # data/fixation_graph_bottleneck/graph/graph_t005/graph_t005_m001
            for transmission_size in args.transmission_size:
                # Create directory
                tsize_path = os.path.join(
                    duration_path, 
                    TSIZE_BASENAME.format(
                        network=network,
                        duration=duration, 
                        tsize=transmission_size,
                    )
                )
                if os.path.exists(tsize_path) and args.overwrite:
                    proc.call(['rm', '-R', tsize_path])
                os.makedirs(tsize_path)

                # data/fixation_graph_bottleneck/graph/graph_t005/graph_t005_m001/graph_t00d_m001_Ns+0
                for Ns in args.Ns:
                    # Create directory
                    ns_path = os.path.join(tsize_path, NS_BASENAME.format(
                        network=network,
                        duration=duration, 
                        tsize=transmission_size, 
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
                            duration=duration, 
                            tsize=transmission_size, 
                            Ns=Ns
                        )
                    )
                    # Overwrite previous log
                    with open(log_path, 'w') as f:
                        print('', file=f)
                    
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
                        generate_single_site_fasta_multihost(
                            pathogen_path, 
                            infected_hostlist=infected_hostlist,
                            npathogens=args.npathogens, 
                            p=args.p
                        )

                        # Write TOML config file inside each instance_path
                        toml_path = os.path.join(
                            instance_path, 
                            TOML_FILENAME.format(
                                network=network,
                                duration=duration, 
                                tsize=transmission_size, 
                                Ns=Ns,
                            )
                        )
                        generate_toml(
                            toml_path, 
                            pathogen_path, 
                            network_path, 
                            fm_path, 
                            instance_path + '/',
                            num_generations=args.generations, 
                            duration=duration, 
                            transmission_prob=args.transmission_prob, 
                            transmission_size=transmission_size,
                            instances=1,
                        )
                        # Run simulation
                        if not args.reduce_output:
                            print('Network:            {}'.format(network))
                            print('Duration:           {}'.format(duration))
                            print('Transmission size:  {}'.format(transmission_size))
                            print('Ns:                 {:+.0f}'.format(Ns))
                            print('Instance:           {}'.format(instance))
                        with open(log_path, 'a') as f:
                            for output in run_contagion(toml_path, threads=args.threads, contagion_path=CONTAGION_PATH):
                                if not args.reduce_output:
                                    print(output.decode('utf-8'), end='')
                                print(output.decode('utf-8'), end='', file=f)

                    # Read log file and write summarary
                    summary_path = os.path.join(
                        basepath, 
                        SUMMARY_FILENAME.format(
                            network=network,
                            duration=duration, 
                            tsize=transmission_size, 
                            Ns=Ns
                        )
                    )
                    summarize_fixation_graph(
                        summary_path, 
                        log_path, 
                        Ns, 
                        network, 
                        reversed_values=True if Ns < 0 else False
                    )

                    # Archive folder
                    archive_path = os.path.join(
                        basepath, 
                        ARCHIVE_FILENAME.format(
                            network=network,
                            duration=duration, 
                            tsize=transmission_size, 
                            Ns=Ns
                        )
                    )
                    proc.call(['tar', '-czf', archive_path, '-C', ns_path, '.'])
                    # Delete folder
                    proc.call(['rm', '-R', ns_path])

    print('Done.')
