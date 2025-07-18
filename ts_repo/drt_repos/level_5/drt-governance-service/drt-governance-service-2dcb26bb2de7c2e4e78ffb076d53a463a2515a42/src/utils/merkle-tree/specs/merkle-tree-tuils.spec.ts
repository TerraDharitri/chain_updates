import { MerkleTreeUtils } from '../markle-tree.utils';
import { promises } from 'fs';

describe('Merkle Tree', () => {
    it('Valid Leaves', () => {
        const leaves = [
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '1',
                votingPower: '1',
            },
            {
                address: 'drt1h2rx7vq29m26ut0rh5x4qnmjk2d93ycukp33jl3at04lhaejzhgqmg50aq',
                balance: '2',
                votingPower: '2',
            },
        ];

        const mp = new MerkleTreeUtils(leaves);

        expect(mp.verifyProof(leaves[0])).toStrictEqual(true);
        expect(mp.verifyProof(leaves[1])).toStrictEqual(true);
    });

    it('Invalid leaves, reversed amounts', () => {
        const leaves = [{
            address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
            balance: '1',
            votingPower: '1',
        }, {
            address: 'drt1h2rx7vq29m26ut0rh5x4qnmjk2d93ycukp33jl3at04lhaejzhgqmg50aq',
            balance: '2',
            votingPower: '2',
        }];

        const badLeaves = [
            {
                address: 'drt1h2rx7vq29m26ut0rh5x4qnmjk2d93ycukp33jl3at04lhaejzhgqmg50aq',
                balance: '1',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '2',
                votingPower: '2',
            },
        ];

        const mp = new MerkleTreeUtils(leaves);

        expect(mp.verifyProof(badLeaves[0])).toStrictEqual(false);
        expect(mp.verifyProof(badLeaves[1])).toStrictEqual(false);
    });

    it('Invalid leaves, wrong addresses', () => {
        const leaves = [{
            address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
            balance: '1',
            votingPower: '1',
        }, {
            address: 'drt1h2rx7vq29m26ut0rh5x4qnmjk2d93ycukp33jl3at04lhaejzhgqmg50aq',
            balance: '2',
            votingPower: '2',
        }];

        const badLeaves = [
            {
                address: 'drt1ha5u92fx54y60lh7wsy7uev8sz2f96gza2n8r0v2ha92jn7ap3ps8rs77l',
                balance: '1',
                votingPower: '1',
            },
            {
                address: 'drt1w547kw69kpd60vlpr9pe0pn9nnqeljrcaz73znenjpgt0h3qlqqqxd8p9v',
                balance: '2',
                votingPower: '2',
            },
        ];

        const mp = new MerkleTreeUtils(leaves);

        expect(mp.verifyProof(badLeaves[0])).toStrictEqual(false);
        expect(mp.verifyProof(badLeaves[1])).toStrictEqual(false);
    });

    it('Depth should be 4', () => {
        const leaves = [
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '1',
                votingPower: '1',
            },
            {
                address: 'drt1h2rx7vq29m26ut0rh5x4qnmjk2d93ycukp33jl3at04lhaejzhgqmg50aq',
                balance: '2',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '3',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '4',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '5',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '6',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '7',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '8',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '9',
                votingPower: '1',
            },
            {
                address: 'drt1sdslvlxvfnnflzj42l8czrcngq3xjjzkjp3rgul4ttk6hntr4qds3x86gw',
                balance: '10',
                votingPower: '1',
            },
        ];

        const mp = new MerkleTreeUtils(leaves);
        expect(mp.getDepth()).toStrictEqual(4);
    });
    it('read from file', async () => {
        const snapshot = `e018d697ad08b3547c49d64e926eed47b0cc5fb56025e8a2941e7f60a4c53fc8`;
        const jsonContent = await promises.readFile(`./src/snapshots/${snapshot}.json`, {
            encoding: 'utf8',
        });
        const leaves = JSON.parse(jsonContent);
        const newMT = new MerkleTreeUtils(leaves);
        expect(newMT.getRootHash()).toEqual(`0x${snapshot}`);
    });
});
