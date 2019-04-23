import {MenuItem} from '../models/menu_item.model';
import {ErcName, InterfaceName, ThemeColor} from './enums';

export const THEME_SETTINGS = {
  color: ThemeColor.LIGHT,
};

export const LOGO_NAMES = {
  [ThemeColor.LIGHT]: 'logo_fullcolor.svg',
  [ThemeColor.DARK]: 'logo_allwhite.svg',
};

export const ROUTES = {
  HOME: 'home',
  BLOCK: 'block',
  ADDRESS_FULL: 'address',
  ADDRESS: 'addr',
  TOKEN: 'token',
  RICHLIST: 'richlist',
  TRANSACTION: 'tx',
  SETTINGS: 'settings',
  VERIFY: 'verify',
  WALLET: 'wallet',
  CREATE_WALLET: 'create-account',
  SEND_TX: 'send-tx',
};

export const MENU_ITEMS: MenuItem[] = [
  {
    title: 'Blocks',
    link: ROUTES.HOME,
    icon: 'fa fa-link fa-fw'
  },
  {
    title: 'Rich List',
    link: ROUTES.RICHLIST,
    icon: 'fa fa-list-ul fa-fw'
  },
  /*{
    title: 'Verify Contract',
    link: '/verify',
    icon: 'fa fa-check-square fa-fw'
  },*/
  {
    title: 'Wallet',
    link: ROUTES.WALLET,
    icon: 'fa fa-wallet fa-fw',
  },
  {
    title: 'Network Stats',
    link: 'https://stats.gochain.io',
    icon: 'fa fa-broadcast-tower fa-fw',
    external: true
  },
  /*{
    title: 'Settings',
    link: '/settings',
    icon: 'fa fa-cogs fa-fw',
  },*/
];

export const DEFAULT_GAS_LIMIT = 21000;

export const TOKEN_TYPES = {
  Erc20: 'ERC20',
  Erc20Burnable: 'ERC20 Burnable',
  Erc20Capped: 'ERC20 Capped',
  Erc20Detailed: 'ERC20 Detailed',
  Erc20Mintable: 'ERC20 Mintable',
  Erc20Pausable: 'ERC20 Pausable',
  Erc165: 'ERC165',
  Erc721: 'ERC721',
  Erc721Receiver: 'ERC721 Receiver',
  Erc721Metadata: 'ERC721 Metadata',
  Erc721Enumerable: 'ERC721 Enumerable',
  Erc820: 'ERC820',
  Erc1155: 'ERC1155',
  Erc1155Receiver: 'ERC1155 Receiver',
  Erc1155Metadata: 'ERC1155 Metadata',
  Erc223: 'ERC223',
  Erc621: 'ERC621',
  Erc777: 'ERC777',
  Erc777Receiver: 'ERC777 Receiver',
  Erc777Sender: 'ERC777 Sender',
  Erc827: 'ERC827',
  Erc884: 'ERC884',
};

export const INTERFACE_ABI = {
  [InterfaceName.GetApproved]: {
    'constant': true,
    'inputs': [
      {
        'name': 'tokenId',
        'type': 'uint256'
      }
    ],
    'name': 'getApproved',
    'outputs': [
      {
        'name': 'operator',
        'type': 'address'
      }
    ],
    'payable': false,
    'stateMutability': 'view',
    'type': 'function'
  },
  [InterfaceName.TransferFrom]: {
    'constant': false,
    'inputs': [
      {
        'name': 'from',
        'type': 'address'
      },
      {
        'name': 'to',
        'type': 'address'
      },
      {
        'name': 'tokenId',
        'type': 'uint256'
      }
    ],
    'name': 'transferFrom',
    'outputs': [],
    'payable': false,
    'stateMutability': 'nonpayable',
    'type': 'function'
  },
  [InterfaceName.SafeTransferFrom]: {
    'constant': false,
    'inputs': [
      {
        'name': 'from',
        'type': 'address'
      },
      {
        'name': 'to',
        'type': 'address'
      },
      {
        'name': 'tokenId',
        'type': 'uint256'
      }
    ],
    'name': 'safeTransferFrom',
    'outputs': [],
    'payable': false,
    'stateMutability': 'nonpayable',
    'type': 'function'
  },
  [InterfaceName.OwnerOf]: {
    'constant': true,
    'inputs': [
      {
        'name': 'tokenId',
        'type': 'uint256'
      }
    ],
    'name': 'ownerOf',
    'outputs': [
      {
        'name': 'owner',
        'type': 'address'
      }
    ],
    'payable': false,
    'stateMutability': 'view',
    'type': 'function'
  },
  [InterfaceName.BalanceOf]: {
    'constant': true,
    'inputs': [
      {
        'name': 'owner',
        'type': 'address'
      }
    ],
    'name': 'balanceOf',
    'outputs': [
      {
        'name': 'balance',
        'type': 'uint256'
      }
    ],
    'payable': false,
    'stateMutability': 'view',
    'type': 'function'
  },
  [InterfaceName.SetApprovalForAll]: {
    'constant': false,
    'inputs': [
      {
        'name': 'operator',
        'type': 'address'
      },
      {
        'name': '_approved',
        'type': 'bool'
      }
    ],
    'name': 'setApprovalForAll',
    'outputs': [],
    'payable': false,
    'stateMutability': 'nonpayable',
    'type': 'function'
  },
  [InterfaceName.SafeTransferFrom1]: {
    'constant': false,
    'inputs': [
      {
        'name': 'from',
        'type': 'address'
      },
      {
        'name': 'to',
        'type': 'address'
      },
      {
        'name': 'tokenId',
        'type': 'uint256'
      },
      {
        'name': 'data',
        'type': 'bytes'
      }
    ],
    'name': 'safeTransferFrom',
    'outputs': [],
    'payable': false,
    'stateMutability': 'nonpayable',
    'type': 'function'
  },
  [InterfaceName.IsApprovedForAll]: {
    'constant': true,
    'inputs': [
      {
        'name': 'owner',
        'type': 'address'
      },
      {
        'name': 'operator',
        'type': 'address'
      }
    ],
    'name': 'isApprovedForAll',
    'outputs': [
      {
        'name': '',
        'type': 'bool'
      }
    ],
    'payable': false,
    'stateMutability': 'view',
    'type': 'function'
  },
  [InterfaceName.Approve]:
    {
      'constant': false,
      'inputs': [
        {
          'name': 'spender',
          'type': 'address'
        },
        {
          'name': 'value',
          'type': 'uint256'
        }
      ],
      'name': 'approve',
      'outputs': [
        {
          'name': '',
          'type': 'bool'
        }
      ],
      'payable': false,
      'stateMutability': 'nonpayable',
      'type': 'function'
    },
  [InterfaceName.TotalSupply]: {
    'constant': true,
    'inputs': [],
    'name': 'totalSupply',
    'outputs': [
      {
        'name': '',
        'type': 'uint256'
      }
    ],
    'payable': false,
    'stateMutability': 'view',
    'type': 'function'
  },
  [InterfaceName.Transfer]: {
    'constant': false,
    'inputs': [
      {
        'name': 'to',
        'type': 'address'
      },
      {
        'name': 'value',
        'type': 'uint256'
      }
    ],
    'name': 'transfer',
    'outputs': [
      {
        'name': '',
        'type': 'bool'
      }
    ],
    'payable': false,
    'stateMutability': 'nonpayable',
    'type': 'function'
  },
  [InterfaceName.Allowance]: {
    'constant': true,
    'inputs': [
      {
        'name': 'owner',
        'type': 'address'
      },
      {
        'name': 'spender',
        'type': 'address'
      }
    ],
    'name': 'allowance',
    'outputs': [
      {
        'name': '',
        'type': 'uint256'
      }
    ],
    'payable': false,
    'stateMutability': 'view',
    'type': 'function'
  },
};

export const ERC_INTERFACE_IDENTIFIERS = {
  [ErcName.Erc20]: [InterfaceName.Allowance, InterfaceName.Approve, InterfaceName.BalanceOf, InterfaceName.TotalSupply, InterfaceName.Transfer, InterfaceName.TransferFrom],
  [ErcName.Erc721]: [InterfaceName.Approve, InterfaceName.BalanceOf, InterfaceName.GetApproved, InterfaceName.IsApprovedForAll, InterfaceName.OwnerOf, InterfaceName.SafeTransferFrom, InterfaceName.SafeTransferFrom1, InterfaceName.SetApprovalForAll, InterfaceName.TransferFrom],
};
