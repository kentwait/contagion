import numpy as np

def generate_random_sequence(chars, num_sites, char_dist=None, seed=None):
    """Returns a list of random string characters.

    Parameters
    ----------
    chars : iterable
        List of single characters to choose from
    num_sites : int
        Length of sequence
    char_dist : list of float or None
        Probability of choosing a character. 
        The position in the list should correspond to the position of the character.
    seed : int or None
        Random seed

    Returns
    -------
    list

    """
    sequence = [''] * num_sites
    if not char_dist:
        char_dist = [1/len(chars)] * len(chars)
    multinomial = np.random.multinomial
    if seed != None:
        # pylint says np.random.RandomState does not exist when it does
        # pylint: disable=E1101
        state = np.random.RandomState(seed)
        multinomial = state.multinomial
    for i, pick in enumerate(multinomial(1, char_dist, size=num_sites)):
        pos = np.argmax(pick)
        sequence[i] = chars[pos]
    return sequence

def construct_fasta_entry(host_id, sequence, description=None, wrap=None):
    """Returns a FASTA-formatted entry for the sequence.

    Parameters
    ----------
    host_id : int
        Host ID of target host
    sequence : iterable
        If sequence is a list of single characters, the list will be joined
        into a string.
    description : str or None
        Description of the sequence. This will be concatenated to the end of
        the description line, ie. >h:1 <description>
    wrap : int or None
        Number of sequence characters per line. If None, the sequence is
        written as a single line.

    Returns
    -------
    str

    """
    wrap = wrap if wrap else len(sequence)
    text = '>h:{host_id}'.format(host_id=host_id)
    if description:
        text += ' {desc}\n'.format(desc=description)
    else:
        text += '\n'
    start = 0
    for end in range(wrap, len(sequence)+1, wrap):
        text += ''.join(sequence[start:end]) + '\n'
        start = end
    return text

def single_host_random_fasta(host_id, n, chars, num_sites, clonal=True, char_dist=None, seed=None, description=None, wrap=None):
    """Returns a FASTA-formatted string of sequences for a single host.

    Parameters
    ----------
    host_id : int
        Host ID of target host.
    n : int
        Number of sequences.
    chars : iterable
        List of single characters to choose from
    num_sites : int
        Length of sequence
    clonal: bool
        Whether to sequences in the host are clones or each sequence 
        is independently randomly generated.
    char_dist : list of float or None
        Probability of choosing a character. 
        The position in the list should correspond to the position of the character.
    seed : int or None
        Random seed
    description : str or None
        Description of the sequence. This will be concatenated to the end of
        the description line, ie. >h:1 <description>
    wrap : int or None
        Number of sequence characters per line. If None, the sequence is
        written as a single line.

    Returns
    -------
    str

    """
    fasta = ''
    seq = ''
    if clonal:
        seq = generate_random_sequence(chars, num_sites, char_dist=char_dist, seed=seed)
    for i in range(n):
        if not clonal:
            seq = generate_random_sequence(chars, num_sites, char_dist=char_dist, seed=seed)
        entry_desc =  '{} count:{}'.format(description, i) if description else 'count:{}'.format(i)
        fasta += construct_fasta_entry(host_id, seq, description=entry_desc, wrap=wrap)
    return fasta

def multiple_host_random_fasta(host_id_list, n, chars, num_sites, clonal='all', 
    char_dist=None, seed=None, description=None, wrap=None):
    """Returns a FASTA-formatted string of sequences for one or more hosts.

    Parameters
    ----------
    host_id_list : list of int
        List of Host ID's of target hosts.
    n : int
        Number of sequences per host.
    chars : iterable
        List of single characters to choose from
    num_sites : int
        Length of sequence
    clonal: 'all', 'host', 'random'
        If 'all', sequences in all hosts are the same. If 'host', sequences in
        each host are the same, but may be different between hosts.
        If 'random', each sequence is independently randomly generated regardless of the host.
    char_dist : list of float or None
        Probability of choosing a character. 
        The position in the list should correspond to the position of the character.
    seed : int or None
        Random seed
    description : str or None
        Description of the sequence. This will be concatenated to the end of
        the description line, ie. >h:1 <description>
    wrap : int or None
        Number of sequence characters per line. If None, the sequence is
        written as a single line.

    Returns
    -------
    str

    """
    assert clonal in ['all', 'host', 'random'], 'clonal value must be one of the following: all, host, random '
    fasta = ''
    seq = ''
    if clonal == 'all':
        seq = generate_random_sequence(chars, num_sites, char_dist=char_dist, seed=seed)
    for host_id in host_id_list:
        if clonal == 'host':
            seq = generate_random_sequence(chars, num_sites, char_dist=char_dist, seed=seed)
        for i in range(n):
            if clonal == 'random':
                seq = generate_random_sequence(chars, num_sites, char_dist=char_dist, seed=seed)
            entry_desc =  '{} count:{}'.format(description, i) if description else 'count:{}'.format(i)
            fasta += construct_fasta_entry(host_id, seq, description=entry_desc, wrap=wrap)
    return fasta
