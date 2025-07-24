import * as crypto from 'crypto-js';
import { BinaryUtils } from "@terradharitri/sdk-nestjs-common";
import '@terradharitri/sdk-nestjs-common/lib/utils/extensions/string.extensions';

export class UsernameUtils {
  private static dnsContracts = [
    'drt1qqqqqqqqqqqqqpgqnhvsujzd95jz6fyv3ldmynlf97tscs9nqqqqgewses',
    'drt1qqqqqqqqqqqqqpgqysmcsfkqed279x6jvs694th4e4v50p4pqqqsmjndc3',
    'drt1qqqqqqqqqqqqqpgqnk5fq8sgg4vc63ffzf7qez550xe2l5jgqqpqygxwm7',
    'drt1qqqqqqqqqqqqqpgqvccl0w78cvr48et3z0n06t8httkp97dkqqps7fjkky',
    'drt1qqqqqqqqqqqqqpgqpzjulx7pemmknndegp2m60vgse3q99zrqqzqz5dm6c',
    'drt1qqqqqqqqqqqqqpgqjr24m738s30aajdv6f5xmzuyxnss8txuqqzsfekdmd',
    'drt1qqqqqqqqqqqqqpgqhcm9k2xkk75e47wpmvfgj8fuzwaguvzyqqrqdhwtys',
    'drt1qqqqqqqqqqqqqpgql2p9sdutmz2vp6w0w5e4f3hjt9wekucvqqrs838yvd',
    'drt1qqqqqqqqqqqqqpgq37fmv57uqkayxctplh4swkw9vfz2jawvqqyqsqlm3z',
    'drt1qqqqqqqqqqqqqpgqlk3vuqlqydznh94ufpt2sy9j4d0u22umqqyscu40uf',
    'drt1qqqqqqqqqqqqqpgqup6wn264pz9c02vklwwrelhtkqydfmh0qq9qu4pz93',
    'drt1qqqqqqqqqqqqqpgq57ekvv62qxhutavp4glfkljkayu8r8rxqq9saqprfc',
    'drt1qqqqqqqqqqqqqpgqcgwaumej906dmftsd8vkyghhpf0ry78eqqxqgg049u',
    'drt1qqqqqqqqqqqqqpgqsf23aanx9pcw60fqmh7vqj2uduzm209aqqxsm4v078',
    'drt1qqqqqqqqqqqqqpgq5yvykcy637980qug55v075ldq3kvntd9qq8qnl2pax',
    'drt1qqqqqqqqqqqqqpgqhygyzhn4ch4lng83cjr4nray80kl5l3sqq8segv3uz',
    'drt1qqqqqqqqqqqqqpgqrdj3kt29t2clcpnnfaevqemf58vtjllqqqgq90cctl',
    'drt1qqqqqqqqqqqqqpgqp64e3pqxwwyy93t5wp2w2jnlf4lfx3ljqqgs2mhd0u',
    'drt1qqqqqqqqqqqqqpgqlmfukwqgzttsld9h8mfm0nqtfu9y8hwjqqfq6ypkst',
    'drt1qqqqqqqqqqqqqpgq3d0uwcsev9dj0d95zvj38xpae7g2az4sqqfsthp0gr',
    'drt1qqqqqqqqqqqqqpgqu3m8g3j684t6f80usm22gfep4x9e5hngqq2q525tp9',
    'drt1qqqqqqqqqqqqqpgqv5r35avg7e94ymw3z2yh8rvmwxdf43clqq2srujwm8',
    'drt1qqqqqqqqqqqqqpgq60584ze6a9udf0q387hzl7t3t87n8q9aqqtqzj399z',
    'drt1qqqqqqqqqqqqqpgqexv860na2t9cwgrvmrydgre23uc5g0ptqqtsghtn6j',
    'drt1qqqqqqqqqqqqqpgqplqkyugu5hesfqmpre80wn4x4zetdmsnqqvqmrqlr0',
    'drt1qqqqqqqqqqqqqpgqrcpy62pjcaxn383dacxt7rmdw3u5am7nqqvs5ve85t',
    'drt1qqqqqqqqqqqqqpgqlfq00tcp9yujyus8540tkk0wjdeg6588qqdq9rvrma',
    'drt1qqqqqqqqqqqqqpgqme5uy9cxp70tehlyzlmsl6kelnh3swqgqqdsgwj86u',
    'drt1qqqqqqqqqqqqqpgqp8w9xf4zjnjfzn6l5pccujc3hvrp2hkcqqwqkdt8qh',
    'drt1qqqqqqqqqqqqqpgq0zq4dx9n5tu65rn39le5tqej88wsgyveqqws8kxkn8',
    'drt1qqqqqqqqqqqqqpgq06hys9ttg4myr7pjy50y03k4uxx387yqqq0qz8vgxy',
    'drt1qqqqqqqqqqqqqpgq5vwmqrwt9r60crhghrxcd2t88hprperlqq0s4ln83z',
    'drt1qqqqqqqqqqqqqpgq62rsz4tq6mp0pjtwsemlxnvd9rakladsqqsqat0gua',
    'drt1qqqqqqqqqqqqqpgqqeu05y7jfjxmc8n34my2nvfjpjm7svxpqqssf0h37h',
    'drt1qqqqqqqqqqqqqpgq98h7d6gnfeg502j4ypp9t5s5gtdczqm4qq3q6jjfhg',
    'drt1qqqqqqqqqqqqqpgq83zzgjxdnt6aucpzmtxnwyqyvc3988f2qq3sncw0pp',
    'drt1qqqqqqqqqqqqqpgqyayxcfwhvv02ettv0ystavn7x3lkuu4tqqjqfghweg',
    'drt1qqqqqqqqqqqqqpgqhmfvs04uzqrjajvslgsypfjhtyyaz7esqqjsuj3ypf',
    'drt1qqqqqqqqqqqqqpgqn948qmzpq8m8lu7y477xh7zlsyndtnm4qqnqlp3ffc',
    'drt1qqqqqqqqqqqqqpgqf9f3kvnr32u6v0nmgch76vnz46ncm7faqqnswrr2l8',
    'drt1qqqqqqqqqqqqqpgqv89a9esd099fvnfj2xyrnnhh584yydejqq5qj35dv2',
    'drt1qqqqqqqqqqqqqpgqmta7xtt292599mray67za5c3rl2yc5h0qq5s5c2yxs',
    'drt1qqqqqqqqqqqqqpgq776u6lt7u5dr6ekn0636t3ua845gfppgqq4qg5wdp4',
    'drt1qqqqqqqqqqqqqpgqy2rhzj2e0nl2mrfxfmgcskmzfl0sdjk2qq4s52pdxc',
    'drt1qqqqqqqqqqqqqpgqpt4q6cny2y9d4qgs4lemkmqy00zwkf9pqqkqr5zmgt',
    'drt1qqqqqqqqqqqqqpgqp05fsn2zvl2tfhlrxz790ss9plkmkzmpqqks8z7ue2',
    'drt1qqqqqqqqqqqqqpgqccz8969nt3d9vw0a64m7zg5px847t0jsqqhqsrkfvm',
    'drt1qqqqqqqqqqqqqpgq64uvd2mq0d07xvt7n8tml6ecnm456qd7qqhs78yf0q',
    'drt1qqqqqqqqqqqqqpgqzq2szpgrh95nj6quqgre2qs3x9hct827qqcqhm3jz5',
    'drt1qqqqqqqqqqqqqpgqsdd0086nxte5h4qrpf9aseexzghkrn4uqqcswrqla2',
    'drt1qqqqqqqqqqqqqpgq4mhmxv3v6878nh3ek8u34kwqqv678fv0qqeqy9gy4l',
    'drt1qqqqqqqqqqqqqpgqa2re7lsvuxz2gcpnkdh7qp75teyqff8gqqesvekht8',
    'drt1qqqqqqqqqqqqqpgqxwljm0xssccjsz4ryr52287rp797eadvqq6qy8had3',
    'drt1qqqqqqqqqqqqqpgqs7kp40xmd4m0xsq8uywtsuxnxuzgltnuqq6smeuncs',
    'drt1qqqqqqqqqqqqqpgq2qg29s5d6vtx7eymm050rhmzetqcqn0dqqmqq76zel',
    'drt1qqqqqqqqqqqqqpgqnwkzwfnm6eqt9rha4fxfz899dv5607cgqqmsyzmm9e',
    'drt1qqqqqqqqqqqqqpgq784s8xvzllndal36c035x7u0wl7cq7phqquq48hpfa',
    'drt1qqqqqqqqqqqqqpgqcuqyx0pcp5fs7qatrzg7cgng2853au6rqqusxra9nf',
    'drt1qqqqqqqqqqqqqpgqxpzxeefc6053qxp3z2sz4jf09ff7cxj8qqaq36glrq',
    'drt1qqqqqqqqqqqqqpgqz35q0ecvzpzhyvum0etghqfcq7d6lsztqqas4pyc7z',
    'drt1qqqqqqqqqqqqqpgqgp42l6n46urv9etd6lxa2ch2m92lclvhqq7qtrg4sy',
    'drt1qqqqqqqqqqqqqpgq2u60t6gppp8uyrtutng27k9quk42xw6qqq7s8dhphj',
    'drt1qqqqqqqqqqqqqpgqm8u4r9hk7ghadf7ufjfq7ywfajxexgf8qqlqg3y5dc',
    'drt1qqqqqqqqqqqqqpgqr098hafftasr030x9jeqeqzchhad99teqqls4jf842',
    'drt1qqqqqqqqqqqqqpgq6ufagmzk0w05hprdasuvgnj5vwfnkjlsqpqqd9xx0h',
    'drt1qqqqqqqqqqqqqpgquf7wpxtnln0a8ywf6hwy82pk5644y2c4qpqs8x83x3',
    'drt1qqqqqqqqqqqqqpgqa76d5tvxqrz9w27ka3wu0ylrdgug3vnaqppq7ygq40',
    'drt1qqqqqqqqqqqqqpgqaq3kkmza02kwj0apclua9kc00e03m0tlqppsgx9na8',
    'drt1qqqqqqqqqqqqqpgqw2p33cncsdtsmqhn8awne3qz37egj85jqpzq0mjq2p',
    'drt1qqqqqqqqqqqqqpgq9qtkqf8sg23lf7zwhyyer2ueypkjsmnpqpzsn4e2zy',
    'drt1qqqqqqqqqqqqqpgq7vl3acd94l2xqw4hglj6lx2ydhq7yy5jqprqcfz83t',
    'drt1qqqqqqqqqqqqqpgqzw0c6t3puxg4rfw9ae8zc2xl3wvmqwxmqprspynt9a',
    'drt1qqqqqqqqqqqqqpgqpwx487sptz3cq3z7jszlz3lkcfdzu397qpyqfyyayx',
    'drt1qqqqqqqqqqqqqpgq3e9e0tqentvwdlratp2m2k9g8n0ljufzqpyss5dc9s',
    'drt1qqqqqqqqqqqqqpgqw82cjja5vqzx635wak76am4y7vu6nnspqp9qvws2mx',
    'drt1qqqqqqqqqqqqqpgqwk47vfp9vme83k34n39nyw9tq3vx9d7eqp9svy3sut',
    'drt1qqqqqqqqqqqqqpgqtnae0xuq9dn54dhphq3y9a7ekzqr7cp7qpxqmepu5m',
    'drt1qqqqqqqqqqqqqpgq58fqfxv7vlhtjlqyl39ufj55zeh4a2jrqpxsaw8yyj',
    'drt1qqqqqqqqqqqqqpgqldp0kd0e0mkxmaay9ntz2j7xytrtm6fpqp8q43mfda',
    'drt1qqqqqqqqqqqqqpgqyaqpjhxfzk06ptke7d4frcy2sc3w456tqp8sc4w5yd',
    'drt1qqqqqqqqqqqqqpgqes2c6e53mu63znddzqxl7az4yrqa3v2eqpgqk85s5f',
    'drt1qqqqqqqqqqqqqpgq3f4s06pp4ra9ysl2yrcnum9654krxhpmqpgs5wdvrr',
    'drt1qqqqqqqqqqqqqpgqa4hmwegr9kp77ywyphps7pu0pzsdq2jxqpfq5g7cc0',
    'drt1qqqqqqqqqqqqqpgqnk3d8yunhr9wndehg3h8y3ku02wgg2deqpfshrymg4',
    'drt1qqqqqqqqqqqqqpgqxmzt96wt6z3pk9tpga8xjgeqld6flq7xqp2qanu7cd',
    'drt1qqqqqqqqqqqqqpgq3ppm8dx4g74zsccz6lp3yyq5pvew59rjqp2s9jlcly',
    'drt1qqqqqqqqqqqqqpgqetuaw54gag7jwtkh8dtdje9j23qc9cqmqptqwuz9uk',
    'drt1qqqqqqqqqqqqqpgqfvl6j0lkppquhx2hhshdjkcd6h2nuflqqpts55qhpy',
    'drt1qqqqqqqqqqqqqpgq83w6606edfk2t0danjsnl64u6tr2ad2sqpvqlrgcjm',
    'drt1qqqqqqqqqqqqqpgqjnq4qedaqd54axzqschax6ullwnk0rn2qpvsu3f097',
    'drt1qqqqqqqqqqqqqpgqykt0f03czqj2p9qltpygzu7jwlzkaxqaqpdqjz074y',
    'drt1qqqqqqqqqqqqqpgqtj4k4u6xxxzqytf2quuy00ggjd0m4e4mqpdszt07kj',
    'drt1qqqqqqqqqqqqqpgqgrvddhuet3vuwnx3vdmshgyve62ar5paqpwq97al56',
    'drt1qqqqqqqqqqqqqpgq7lrcmp7uhlxxp6jwa7s8du6zkyhl7sp0qpws3zl2aq',
    'drt1qqqqqqqqqqqqqpgqus7asxmg46uasmcmrwvyjgyhqtpusua4qp0qqdh0px',
    'drt1qqqqqqqqqqqqqpgqkvpyh4pw97asqrgmpnhvedtwurc5tkyvqp0s7e6wcf',
    'drt1qqqqqqqqqqqqqpgq5mk4uucfw9h6y0swehrwdmachxvn4hu0qpsqexttml',
    'drt1qqqqqqqqqqqqqpgqvqqpd2fyqf7ypre892smscarfagast3gqpsscuktz9',
    'drt1qqqqqqqqqqqqqpgqur88h5n8t8whtfurlfta4x0uycppg392qp3qlxaycu',
    'drt1qqqqqqqqqqqqqpgq94atq5qqluaeuewx9mdj0klhedjqmrjvqp3shcldfq',
    'drt1qqqqqqqqqqqqqpgqna20c6fhfdl7tfrkwjtnhtxayajrckasqpjq3x7g57',
    'drt1qqqqqqqqqqqqqpgq7qs0ch9c7mphjecr59rwtk75ptss2gw6qpjsdheuqp',
    'drt1qqqqqqqqqqqqqpgqz0ycyug2rqtpyrh5p33y9vqjv95s3xmaqpnqrq4jr7',
    'drt1qqqqqqqqqqqqqpgq3qd7r8025xqxqh9vzygmllwecuxyrq4cqpnsslswxw',
    'drt1qqqqqqqqqqqqqpgqk6lw9w7aj5scgqh23f0juk55ydan0rzmqp5qdjv5uu',
    'drt1qqqqqqqqqqqqqpgqywkkwsumtqvuncmx3exjrm9vqp76u3vyqp5sszf56x',
    'drt1qqqqqqqqqqqqqpgqydsrr9rw6wgl83dm7mgyqzr7khqp982zqp4q36cj29',
    'drt1qqqqqqqqqqqqqpgq058jfs3fnd6tjdpg0gfdg5pd57pp672nqp4sf4vylz',
    'drt1qqqqqqqqqqqqqpgqe4crz3fxffsag57ze6jnyn70r0zrkz6qqpkqpfenph',
    'drt1qqqqqqqqqqqqqpgqytqyk2zxplmr0znx8q0tkwrwf6pnn5muqpksshhhtd',
    'drt1qqqqqqqqqqqqqpgqr9nlxmhpfhasz9y7504mtp6pdqmh98p8qphq4azy9c',
    'drt1qqqqqqqqqqqqqpgqe2cmllq3zhwfuzdpdzqh7223xnc907ffqphs6xrqlh',
    'drt1qqqqqqqqqqqqqpgqawhqqvutzqmz2n7sn2tame55yaejx00vqpcq8upk7p',
    'drt1qqqqqqqqqqqqqpgql0x09mnzwn4hwkfpsytn9uhjyz0t5g5rqpcs7uz84h',
    'drt1qqqqqqqqqqqqqpgqf4dpre4qv5fca6a4y3vez5rznpfh4k82qpeqm05sgw',
    'drt1qqqqqqqqqqqqqpgq6q3jfw7l7fw7vd02de2nn3rgxvnv29lkqpessc4r23',
    'drt1qqqqqqqqqqqqqpgqz7s9ut07qm9uh2rtzzuntnnuxrpejjscqp6qjder0t',
    'drt1qqqqqqqqqqqqqpgq7yqrph2w055swrjjellgu55r4vjnfyklqp6st8gt24',
    'drt1qqqqqqqqqqqqqpgq9nq5dgclefa6l9np6pvtgn4mpmh8zt4nqpmq95zhj8',
    'drt1qqqqqqqqqqqqqpgqfhc0kxmy8ryn5ezew0fwx7ulnfggjtzpqpmsa747fe',
    'drt1qqqqqqqqqqqqqpgqrqktg0vqjkwccz47gl5zcnvnqm6tq0srqpuqsft68p',
    'drt1qqqqqqqqqqqqqpgq7rxe4kmjdxpj4cfnlh96xeqva9yjx5jhqpus09czek',
    'drt1qqqqqqqqqqqqqpgq6amzugw7tqlqgwamdel04m437pt9rldpqpaqwwaxk7',
    'drt1qqqqqqqqqqqqqpgqjverqd04c8ya7fa4aumm26rxka59tvxmqpasvmwc6y',
    'drt1qqqqqqqqqqqqqpgqxnkx4p5u8gzxwd2ekqqsa7aktr6trphdqp7qxl97w3',
    'drt1qqqqqqqqqqqqqpgq0474ralfs6086l0a0zqz0dtf3clnpe2lqp7slxykdh',
    'drt1qqqqqqqqqqqqqpgqangd34sj05xsss6yp84l8j8hteqmqp36qplqcdk6fu',
    'drt1qqqqqqqqqqqqqpgqsy2c35ne2a076z3cudpzgdjm2yp2uuhjqpls54hdha',
    'drt1qqqqqqqqqqqqqpgqg9ugsf9g7vcyh6ahyfshjqjyr65y5sw7qzqq0sj0nm',
    'drt1qqqqqqqqqqqqqpgqmnk6qv550mycf7fpgvk2yfr0e8vy7jhwqzqskmarw8',
    'drt1qqqqqqqqqqqqqpgqjwltr4xvafv2edc958klff3u0qlnem94qzpq79480p',
    'drt1qqqqqqqqqqqqqpgqz2aa6uaptw8ryf254u93xqwl4lkvu9vmqzpsmvqlhq',
    'drt1qqqqqqqqqqqqqpgq67c6t69pcc7dt597xmlnwjr4c55ws55qqzzqqgx84x',
    'drt1qqqqqqqqqqqqqpgqf97pgqdy0tstwauxu09kszz020hp5kgqqzzsdyhgds',
    'drt1qqqqqqqqqqqqqpgqdlfj9r696jpmsakvqp7redflu472c3zgqzrq542zyq',
    'drt1qqqqqqqqqqqqqpgqdyp9pfyj9ueurwxd8p42na58vsn5r9pzqzrsqsnvt2',
    'drt1qqqqqqqqqqqqqpgqtprnv8gasxjkt79jn5aq9xrantj9laa6qzyqrkdctj',
    'drt1qqqqqqqqqqqqqpgq0wa5j9l6urteh5k0agkkg9af364y77x5qzysd0u4nk',
    'drt1qqqqqqqqqqqqqpgq0k3kl9kx9swl9klueljdplv8x9m9zyf4qz9q6xce8k',
    'drt1qqqqqqqqqqqqqpgq6g34lngl7w9w4yvq0fr3u72xushkssusqz9spzsjdx',
    'drt1qqqqqqqqqqqqqpgqf4pw79l5s9xkslyf5p06egwcnjul95ksqzxqvw9760',
    'drt1qqqqqqqqqqqqqpgqkvve4g87rg628mvfmh0rdwe2kkqqskc4qzxsjvhz8z',
    'drt1qqqqqqqqqqqqqpgqycngxuhcjcw0wv5gpe58h0pw6g89cthdqz8q4y98wl',
    'drt1qqqqqqqqqqqqqpgq2dxd9gfff6u4pv8w9ry89cuksy34kjkcqz8sudma7t',
    'drt1qqqqqqqqqqqqqpgq2leexk6fwaxlxggzrnkxzruwsjzfcq2mqzgq7hwwcs',
    'drt1qqqqqqqqqqqqqpgqg0vmcgrmrqnkec0fem3rf2l8468th3znqzgstuqsxv',
    'drt1qqqqqqqqqqqqqpgqefmp3c3pzwetn9xt5eav7ttxsnpg44aqqzfqp7luf5',
    'drt1qqqqqqqqqqqqqpgqs0cxu3qq405r7ua8xheys8nkmz25sar3qzfswgcvad',
    'drt1qqqqqqqqqqqqqpgq8v7py8qke4kd6l3np0a2hfs3zz637782qz2qpulcg7',
    'drt1qqqqqqqqqqqqqpgq6j24n4kyyk77rdmduzzv3x5v74hg3dqnqz2sw0ys3l',
    'drt1qqqqqqqqqqqqqpgq3a2fedxlf9n9ga5qt9znj5v6eyuunqm7qztq7ta4mj',
    'drt1qqqqqqqqqqqqqpgq2jq3q6k9p007v8n58k2qqetfck6pt0pcqztsj5lp00',
    'drt1qqqqqqqqqqqqqpgq8ega5ywe9y6l2yt664cssdmeuanhnxs4qzvq0laq5h',
    'drt1qqqqqqqqqqqqqpgqv0vvklxxzvxyq9nx0zvnlkus2a3c8267qzvs6h4t7z',
    'drt1qqqqqqqqqqqqqpgqwk2a4cg4vgx0m2ummjf8x885a4y4ypdeqzdqu4vm3l',
    'drt1qqqqqqqqqqqqqpgqrmwah0xxzax9g2geyk99twlhe4yfrkgmqzdsmsa59g',
    'drt1qqqqqqqqqqqqqpgqynp5z59pgqjnphfxg00mkpzx54an6038qzwqrpkmd4',
    'drt1qqqqqqqqqqqqqpgqvney9xtjnp4duxl5y2slmg3k72g7vj5vqzws3kp5a6',
    'drt1qqqqqqqqqqqqqpgqrnkmx4us6qw3wprmlmltuyqqvzy6zwkfqz0qrx7mxp',
    'drt1qqqqqqqqqqqqqpgqx4ca3eu4k6w63hl8pjjyq2cp7ul7a4ukqz0stud29v',
    'drt1qqqqqqqqqqqqqpgqpmv02e6a5k7elq0qwp77t9cynhq4tt0aqzsqvt4cym',
    'drt1qqqqqqqqqqqqqpgq7fjjs979sj4pxpeuru34ku85n2j3lj2aqzss5x7gwm',
    'drt1qqqqqqqqqqqqqpgqcp8f3lv7kw4ckhq6y3yfv7awh45gf5hrqz3q5xpd57',
    'drt1qqqqqqqqqqqqqpgqc6nz2zmkxnlfwjzg9jj8275rwdxzv859qz3sxedv0j',
    'drt1qqqqqqqqqqqqqpgqcl8p0wzfluw2vvwnud57hn7drt0mma6hqzjqj7y2ja',
    'drt1qqqqqqqqqqqqqpgq2yue5t96xetjkt0rfvjn5waf2zshztx8qzjswpg58s',
    'drt1qqqqqqqqqqqqqpgq46xccpqq3vd5e9ffyt4804e9ln5rxullqznq99lsfe',
    'drt1qqqqqqqqqqqqqpgqvry6tjmaa60w0k4va85a927l20k3hmmxqznswgahjl',
    'drt1qqqqqqqqqqqqqpgqze3nszzj6ngfw3g0thm30r3v4tzudexdqz5q824gc2',
    'drt1qqqqqqqqqqqqqpgqxhe0c0d7lguv3hn267sa9zkdwrz0hc5vqz5sjhhnd9',
    'drt1qqqqqqqqqqqqqpgq3lv2y7rnuvd2rux5d00cs2plhyl6vpsfqz4qz3ujeh',
    'drt1qqqqqqqqqqqqqpgqvrsdh798pvd4x09x0argyscxc9h7lzfhqz4skhg7gk',
    'drt1qqqqqqqqqqqqqpgqkfpqz9ypt9rwxr42ygrgjg0wluep49t8qzkqe0nhj7',
    'drt1qqqqqqqqqqqqqpgqpvx058h36yczqsdmj4cyfn37cy0gxt43qzksa2svn7',
    'drt1qqqqqqqqqqqqqpgqrx5fz9lsd0nz526wmcjnj5cv6as3y2qkqzhqfs6qxr',
    'drt1qqqqqqqqqqqqqpgqqa39pkddyw4plx78nw686qgxatd7ukthqzhskwve2u',
    'drt1qqqqqqqqqqqqqpgqjp8ev4hc4j3sslmqg0kyeyhtlh4vm9fpqzcqdlfsrj',
    'drt1qqqqqqqqqqqqqpgq5xlqkfxcf4rpta6mkwrufgfpgw6r3nf7qzcsjvfdlj',
    'drt1qqqqqqqqqqqqqpgqfq9na9k5jh7zpvly262z3k2yrmvdlpvfqzeq8570r6',
    'drt1qqqqqqqqqqqqqpgqvpc85zpq9e5edm4velyuxqejmt49ldfkqzesqws74j',
    'drt1qqqqqqqqqqqqqpgqmsqw7qx5hpnvrkc6vummdzf60e93dswrqz6q0p7257',
    'drt1qqqqqqqqqqqqqpgq0hsaa6whyjclpfqnyvuaqkgaxwt9ncruqz6s0yjatv',
    'drt1qqqqqqqqqqqqqpgqa7szyvngye2j78axcdcj534hp3wjprauqzmqlux74v',
    'drt1qqqqqqqqqqqqqpgqrmeqdfgucxzx6jekzhdrhly8789hjsx6qzmsqmsuex',
    'drt1qqqqqqqqqqqqqpgq4su5u6t3khfkd6u853cudq2xqtzrksspqzuqnnk250',
    'drt1qqqqqqqqqqqqqpgq4tq46jzq2gfwrglvv8a6fgr44vu87wjsqzusszaytu',
    'drt1qqqqqqqqqqqqqpgqe503urc4w4t4900h7qlf4m89xn0zxed7qzaqtkd29w',
    'drt1qqqqqqqqqqqqqpgqyykmg4r0zk7j8e2vgd2u8nl57up6fkpsqzaszlx25f',
    'drt1qqqqqqqqqqqqqpgqexzt863hvc786ccjg6qz8qmqlw5slkhhqz7q7rlyhv',
    'drt1qqqqqqqqqqqqqpgqp0qwk8xc53y6c7ae3uj3msg57c65n7c4qz7s7kyzha',
    'drt1qqqqqqqqqqqqqpgqr00f20jkjzj6jdwyyfuhs6hz8dj2c2jgqzlqwkapj4',
    'drt1qqqqqqqqqqqqqpgqapjemgty9de39pf6760694gz52ekusgmqzls67q5jr',
    'drt1qqqqqqqqqqqqqpgql47hud42jcskpggf3xzxpputfc6xprmfqrqqqesr3s',
    'drt1qqqqqqqqqqqqqpgqs44hhz7z32mjzj0sn00glncl2pf5p4kkqrqsdvv4m9',
    'drt1qqqqqqqqqqqqqpgq0fun5y08yda48a6avywsds86t9cy656hqrpq5a39hj',
    'drt1qqqqqqqqqqqqqpgqdc50p7fxtgkwdwal2xnrk2f5df97wyrdqrpszs52wp',
    'drt1qqqqqqqqqqqqqpgqt8t6hj0mxxmwj6zrwlldz9n7pw9eutuaqrzquhv4zd',
    'drt1qqqqqqqqqqqqqpgq42cn7rfgjd7gq859thd5x8sehdm7c6nwqrzsft08gl',
    'drt1qqqqqqqqqqqqqpgqfuh0mnpdf4ej2pfdd63lt46gyny0senfqrrqy26nrl',
    'drt1qqqqqqqqqqqqqpgqypq2ux9gdyvw4lkfy57sgprzeaxyu3saqrrs7wft9v',
    'drt1qqqqqqqqqqqqqpgq5nqp808kuwpqfkyqdkns29lv0rhxdz6pqryqus0fa9',
    'drt1qqqqqqqqqqqqqpgq7e03tpmhuqhacqdfaukn50wkllxpmzetqrysmgafx8',
    'drt1qqqqqqqqqqqqqpgqfnsltcmhdlyyjqlhv8sgr0p3rzzs839mqr9qqqhm5w',
    'drt1qqqqqqqqqqqqqpgqgpdjf5ptrge0leahucfhagvkvmxnknfkqr9snndtc3',
    'drt1qqqqqqqqqqqqqpgqjw7ytmn5jp28zngmsq9gdsgslp49s3wdqrxq6wgftn',
    'drt1qqqqqqqqqqqqqpgqg3nvuz35s762vea4gyww5nn0h87r3987qrxs794gfx',
    'drt1qqqqqqqqqqqqqpgqhjhhs4920jfxkastaf9ze8yyd4jxws6wqr8qlxuunq',
    'drt1qqqqqqqqqqqqqpgqdpkgt7da7xdpfzc99lwd8fcp9aqcem3nqr8skjeg9n',
    'drt1qqqqqqqqqqqqqpgq5pua3suerg6c8de4vucexrq2c0d5xctjqrgq9rrhjf',
    'drt1qqqqqqqqqqqqqpgqyt7avg7te9mngysa07dw04d49agzd9deqrgs4ugxs7',
    'drt1qqqqqqqqqqqqqpgqvw7yxu6055ucua2ytnfd04drdju2aun9qrfq3fjdgn',
    'drt1qqqqqqqqqqqqqpgq4aae5mpfxnywf5yapjm2c2nudgg5ewl3qrfsvwuxcn',
    'drt1qqqqqqqqqqqqqpgq9ddcdjrq55w6kp8pvfa8d7m8dsyrjzjuqr2q4jz3j3',
    'drt1qqqqqqqqqqqqqpgqs3kh4kjyckz5zwqzy54spxx27zss6rvvqr2sce2efq',
    'drt1qqqqqqqqqqqqqpgqfka5dqs398kumyqfvcpthd0f43x6htq7qrtqj52w57',
    'drt1qqqqqqqqqqqqqpgqdsg4wn28q8j8w5ct5s2cxr04spjul8cuqrtshgqxgf',
    'drt1qqqqqqqqqqqqqpgq9g7umcgkf8jcvv6u93g8cl78znt4afa9qrvq0779af',
    'drt1qqqqqqqqqqqqqpgqzsqek20mul9t07kwj8fswjl85qnmysweqrvsgqdrr0',
    'drt1qqqqqqqqqqqqqpgqe5f8pf52umknu489z4tuat3h4tql3e9mqrdqye6quf',
    'drt1qqqqqqqqqqqqqpgqnp5x7hvss26dufx5p0csls30g7wuquhxqrdse9wzku',
    'drt1qqqqqqqqqqqqqpgqhh8w5q9mvv9tn6w9e768htmj4t9y4390qrwqj76pgg',
    'drt1qqqqqqqqqqqqqpgqx4m8kz3jmespyr47tyrhgrmzsv9qlk32qrws9p92lm',
    'drt1qqqqqqqqqqqqqpgqz5gjr6e8c2qcpqawepr0t2urnkq07p4aqr0qhje5tn',
    'drt1qqqqqqqqqqqqqpgq8pa6mkfyt8lvq6cl3ppgnskn9ms03ne3qr0smh8yy7',
    'drt1qqqqqqqqqqqqqpgqjrzh4wfkkkmkk7ke3upqnmmr3amq7va2qrsqus7y7l',
    'drt1qqqqqqqqqqqqqpgqfkslzjt70pz7py3vjemvmshewgemvglmqrss8r37tu',
    'drt1qqqqqqqqqqqqqpgqyykpxr693gc3mpvvc4k3frq8m7eeguuxqr3qcm26mt',
    'drt1qqqqqqqqqqqqqpgqq68jetqwtapzmstwnz7ex8s9tqy25kheqr3scg90ze',
    'drt1qqqqqqqqqqqqqpgqr3pqqnz626ja7sp2j35mjq9qwt25lyf3qrjqxu6yc4',
    'drt1qqqqqqqqqqqqqpgqqdkg0chrthsjm30u4rlw8etklrp3qx3aqrjs8zw4fy',
    'drt1qqqqqqqqqqqqqpgqcqc0dls8kdykj2gk2ckcdkrehf5m7wrrqrnq7fl235',
    'drt1qqqqqqqqqqqqqpgq29pf3z3rt23smjdlppnp0w65sq0hyh92qrnsrsh2cm',
    'drt1qqqqqqqqqqqqqpgqm0gn9hh52s9qqusxxye2rxlqcphfscuwqr5q84kx77',
    'drt1qqqqqqqqqqqqqpgq5wqgz3tkrr90pdq07sh2qkgvxdujpsmwqr5s37n0fd',
    'drt1qqqqqqqqqqqqqpgqp58nd7xs43zqr43s2aq2z3km4mqmgjvjqr4q29ps5m',
    'drt1qqqqqqqqqqqqqpgq8wrk96xm882m478srwkwlnr9ugjfvhqwqr4sqjh3zm',
    'drt1qqqqqqqqqqqqqpgq4wek8nyjsgd6jr3ctj0zzkgvjts8me6wqrkq5q24tp',
    'drt1qqqqqqqqqqqqqpgqlyfmgtd30x36zdpuxydvpkydtvupf2hsqrksyrlnm4',
    'drt1qqqqqqqqqqqqqpgq45p6f3pem37kjgkh238a0urmayj7k49dqrhqvmw5nv',
    'drt1qqqqqqqqqqqqqpgqc0htpys8vhtf5m3tg7t6ts2wvkgx3favqrhssv4xd3',
    'drt1qqqqqqqqqqqqqpgqrm4lhgjag99lm5k0nrqf0fhu9avl56n2qrcqfvyn8e',
    'drt1qqqqqqqqqqqqqpgqsuxg5m7p4ka5d5q8x2t9x6nxm5356chfqrcsclzc3j',
    'drt1qqqqqqqqqqqqqpgqez3jrxx3krhncmzhmzf8ksnjvmt0x6vrqreq39wnyl',
    'drt1qqqqqqqqqqqqqpgq05sran5e9w3kwzu4mkz7prtlvqeez583qres3sa5d2',
    'drt1qqqqqqqqqqqqqpgq3m44gud3g9z5p7xjevm80dk2g504fhzhqr6qma2cx6',
    'drt1qqqqqqqqqqqqqpgqrequgw4plxarc0rcu49gtxc7874g7waxqr6srn0gls',
    'drt1qqqqqqqqqqqqqpgqtxf5ra4qkj4cmgv8q070w8g2cyd8y56jqrmqc6yw80',
    'drt1qqqqqqqqqqqqqpgq3uxwmwtgmms6jytn3vzlw89vrxxe9xjwqrms093fy0',
    'drt1qqqqqqqqqqqqqpgqs6d64ffr3m4veck8cuvkp8zgn2gflmemqruqnlq9nh',
    'drt1qqqqqqqqqqqqqpgqt9moaedq83vc8fdqjfknuend4auj7gewqrusp6qvqd',
    'drt1qqqqqqqqqqqqqpgqwq64a8w0gqqnxchm572tj38uda9elqw7qraq6m8a34',
    'drt1qqqqqqqqqqqqqpgqzwl6axaq2kfjaglfsrm3eqc0rs3cn85rqraskkfe8m',
    'drt1qqqqqqqqqqqqqpgq3e4gtyegzqttlfqau6wtzm7vg4hl2fjmqr7qqgky0n',
    'drt1qqqqqqqqqqqqqpgqf45egn6zvc39q35m5ffw8x4h8c9p2929qr7sdxllch',
    'drt1qqqqqqqqqqqqqpgqlkk3dj24u5gfn75ymsd3g490qya7xzumqrlqjdwwel',
    'drt1qqqqqqqqqqqqqpgqd3z70c44v7g9makg800l3r3gx3e5a97xqrlszuer7m',
  ];

  static normalizeUsername(username: string): string {
    const prefix = '@';
    const suffix = '.numbat';

    let normalizedUsername = username.toLowerCase().replace(/\W/g, '');

    if (normalizedUsername.startsWith(prefix)) {
      normalizedUsername = normalizedUsername.substring(prefix.length);
    }

    if (!normalizedUsername.endsWith(suffix)) {
      normalizedUsername += suffix;
    }

    return normalizedUsername;
  }

  static getContractAddress(username: string): string {
    const normalized = UsernameUtils.normalizeUsername(username);

    const hash = crypto.SHA3(normalized, { outputLength: 256 }).toString(crypto.enc.Hex);

    const buffer = Buffer.from(hash, 'hex');
    const last = buffer[buffer.length - 1];

    return UsernameUtils.dnsContracts[last];
  }

  static extractUsernameFromRawBase64(rawUsername: string): string {
    if (!rawUsername || rawUsername.length == 0 || !this.isBase64(rawUsername)) {
      return '';
    }

    const decodedUsername = BinaryUtils.base64Decode(rawUsername);
    if (decodedUsername.length == 0) {
      return '';
    }

    return decodedUsername;
  }

  // TODO: might move this to BinaryUtils
  private static isBase64(str: string): boolean {
    if (!str) {
      return false;
    }

    return str.length % 4 == 0 && /^[A-Za-z0-9+/]+={0,2}$/.test(str);
  }

  static encodeUsername(username: string): string {
    const normalized = UsernameUtils.normalizeUsername(username);
    return Buffer.from(normalized).toString('hex');
  }
}
