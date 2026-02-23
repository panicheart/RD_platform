-- Seed data for product shelf
-- Creates 3 example products

-- Insert example products
INSERT INTO products (id, name, description, trl_level, category, version, is_published, download_count, metadata, created_at, updated_at)
VALUES 
(
    '01HQ123456789012345678901',
    '宽带射频接收通道模块',
    '适用于2-18GHz频段的宽带射频接收通道，包含低噪声放大器、混频器、滤波器等组件。适用于A/B/C类产品的通用接收前端。',
    7,
    '射频组件',
    'v2.1.0',
    true,
    15,
    '{"frequency_range": "2-18GHz", "noise_figure": "<3dB", "gain": "30dB", "power": "5V/200mA"}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    '01HQ123456789012345678902',
    '软件定义无线电(SDR)基带处理板',
    '基于Zynq UltraScale+ RFSoC的SDR基带处理板，集成ADC/DAC，支持100MHz瞬时带宽，可用于多种波形体制的数字信号处理。',
    8,
    '数字组件',
    'v1.5.2',
    true,
    23,
    '{"adc_resolution": "14bit", "dac_resolution": "14bit", "bandwidth": "100MHz", "processing": "ZU47DR"}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    '01HQ123456789012345678903',
    '小型化频率合成器',
    '低相位噪声频率合成器模块，输出频率范围100MHz-6GHz，相位噪声优于-110dBc/Hz@10kHz，适用于频率源、本振等场景。',
    9,
    '频率组件',
    'v3.0.1',
    true,
    42,
    '{"frequency_range": "100MHz-6GHz", "phase_noise": "-110dBc/Hz@10kHz", "switching_speed": "<10us", "size": "50x30mm"}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Add versions for first product
INSERT INTO product_versions (id, product_id, version, changelog, created_at)
VALUES 
(
    '01HQ987654321098765432109',
    '01HQ123456789012345678901',
    'v1.0.0',
    '初始版本，完成基本功能验证',
    CURRENT_TIMESTAMP - INTERVAL '6 months'
),
(
    '01HQ98765432109876543210A',
    '01HQ123456789012345678901',
    'v2.0.0',
    '优化噪声系数，扩展频段至2-18GHz',
    CURRENT_TIMESTAMP - INTERVAL '3 months'
),
(
    '01HQ98765432109876543210B',
    '01HQ123456789012345678901',
    'v2.1.0',
    '改进电源滤波，提升抗干扰能力',
    CURRENT_TIMESTAMP
);

-- Add versions for second product
INSERT INTO product_versions (id, product_id, version, changelog, created_at)
VALUES 
(
    '01HQ98765432109876543210C',
    '01HQ123456789012345678902',
    'v1.0.0',
    '初始版本，基于ZU27DR',
    CURRENT_TIMESTAMP - INTERVAL '8 months'
),
(
    '01HQ98765432109876543210D',
    '01HQ123456789012345678902',
    'v1.5.0',
    '升级到ZU47DR，带宽提升至100MHz',
    CURRENT_TIMESTAMP - INTERVAL '2 months'
),
(
    '01HQ98765432109876543210E',
    '01HQ123456789012345678902',
    'v1.5.2',
    '修复温度补偿问题',
    CURRENT_TIMESTAMP
);

-- Add versions for third product
INSERT INTO product_versions (id, product_id, version, changelog, created_at)
VALUES 
(
    '01HQ98765432109876543210F',
    '01HQ123456789012345678903',
    'v1.0.0',
    '初始版本，基于ADF4351',
    CURRENT_TIMESTAMP - INTERVAL '2 years'
),
(
    '01HQ987654321098765432110',
    '01HQ123456789012345678903',
    'v2.0.0',
    '改用集成VCO方案，体积缩小50%',
    CURRENT_TIMESTAMP - INTERVAL '1 year'
),
(
    '01HQ987654321098765432111',
    '01HQ123456789012345678903',
    'v3.0.0',
    '全新架构，相位噪声优化10dB',
    CURRENT_TIMESTAMP - INTERVAL '6 months'
),
(
    '01HQ987654321098765432112',
    '01HQ123456789012345678903',
    'v3.0.1',
    '增加SPI配置接口',
    CURRENT_TIMESTAMP
);
