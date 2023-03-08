# Remove previous config data
rm -rf ~/.anchor
make install

# Anchor initialize
anc i

# Recover the key in the main chain 
anc k recover

# Send transaction to main chain in order to store the anchor contract
anc e ctrt store

# Send transaction to main chain in order to instantiate the anchor contract
anc e ctrt inst

# Start the anchor
anc e s